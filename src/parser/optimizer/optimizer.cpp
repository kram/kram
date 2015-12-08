// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "optimizer.h"

#include <unordered_map>
#include <string>
#include <iostream>

std::vector<Instruction*> Optimizer::variable_alloc(std::vector<Instruction*> instructions)
{
	auto opt = new Optimizer();
	return opt->variable_alloc_level(instructions);
}

Optimizer::Optimizer()
{
	this->depth = 0;
	this->next_num = {1};
	this->names_map = {std::unordered_map<std::string, size_t>{}};
}

void Optimizer::push()
{
	auto m = std::unordered_map<std::string, size_t>{};

	this->names_map.push_back(m);
	this->next_num.push_back(1);

	this->depth++;
}

void Optimizer::pop()
{
	this->names_map.pop_back();
	this->next_num.pop_back();

	this->depth--;
}

stack_and_pos Optimizer::variable_alloc_resolve_name(std::string name)
{
	auto i = this->depth;

	while (true) {

		if (this->names_map[i].count(name) == 1) {
			return stack_and_pos{i, this->names_map[i][name]};
		}

		if (i == 0) {
			return stack_and_pos{0, 0};
		}

		--i;
	}

	return stack_and_pos{0, 0};
}

std::vector<Instruction*> Optimizer::variable_alloc_level(std::vector<Instruction*> instructions)
{
	for (Instruction* ins : instructions) {

		bool did_push_stack = false;

		if (ins->instruction == Ins::FUNCTION ||
			ins->instruction == Ins::WHILE ||
			ins->instruction == Ins::IF ||
			ins->instruction == Ins::DEFINE_CLASS)
		{
			this->push();
			ins->stack_and_pos = stack_and_pos{this->depth, 0};
			did_push_stack = true;
		}

		if (ins->instruction == Ins::ASSIGN) {
			ins->stack_and_pos = this->variable_alloc_resolve_name(ins->name);

			if (KR_SNP_GET_POS(ins->stack_and_pos) == 0) {
				ins->stack_and_pos = stack_and_pos{this->depth, this->next_num[this->depth]};

				// std::cout << "Defined " << ins->name << " as " << KR_SNP_GET_STACK(ins->stack_and_pos) << ":" << KR_SNP_GET_POS(ins->stack_and_pos) << "\n";

				this->names_map[this->depth][ins->name] = this->next_num[this->depth];
				++this->next_num[this->depth];
			}
		}

		if (ins->instruction == Ins::FUNCTION) {
			this->names_map[this->depth]["self"] = this->next_num[this->depth];
			++this->next_num[this->depth];
		}

		if (ins->instruction == Ins::FUNCTION_PARAMETER) {
			ins->stack_and_pos = stack_and_pos{this->depth, this->next_num[this->depth]};

			this->names_map[this->depth][ins->name] = this->next_num[this->depth];
			++this->next_num[this->depth];
		}

		if (ins->instruction == Ins::NAME || ins->instruction == Ins::SET) {
			ins->stack_and_pos = this->variable_alloc_resolve_name(ins->name);

			if (KR_SNP_GET_POS(ins->stack_and_pos) == 0) {
				//std::cout << "Optimizer could not find name...? " << ins->name << "\n";
			}
		}

		if (ins->left.size() > 0) {
			this->variable_alloc_level(ins->left);
		}

		if (ins->center.size() > 0) {
			this->variable_alloc_level(ins->center);
		}

		if (ins->right.size() > 0) {
			this->variable_alloc_level(ins->right);
		}

		if (did_push_stack) {
			this->pop();
		}

	}

	return instructions;
}