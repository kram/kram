// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

#include <iostream>

size_t VM::env_get_pos(stack_and_pos snp)
{
	return this->env_get_pos(KR_SNP_GET_STACK(snp), KR_SNP_GET_POS(snp));
}

size_t VM::env_get_pos(size_t stack, size_t pos)
{
	return this->env_stack_positions[stack].back() + pos;
}

void VM::name_create(const std::string& name, Value* val)
{
	this->environment->set(name, val);
}

void VM::name_create(size_t pos, Value* val)
{
	if (pos > this->env_depth_current_max) {
		this->env_depth_current_max = pos;
	}

	this->environment->set(pos, val);
}

void VM::name_create_root(std::string name, Value* val)
{
	this->environment->set_root(name, val);
}

void VM::name_update(size_t pos, Value* val)
{
	this->environment->update(pos, val);
}

void VM::name_update(const std::string& name, Value* val)
{
	// Verify that the variable exists first
	if (!this->environment->exists(name)) {
		std::cout << "No such variable, " << name << ", did you mean to use := ?\n";
		exit(0);
	}

	Value* previous = this->environment->get(name);

	if (previous->type != val->type) {
		std::cout << "Can not update (with =) a variable of type " << previous->print(true) << " to " << val->print(true) << "\n";
		exit(0);
	}

	this->environment->update(name, val);
}

Value* VM::name_get(const std::string& name)
{
	return this->environment->get(name);
}

Value* VM::name_get(size_t pos)
{
	return this->environment->get(pos);
}

Value* VM::name_get_root(const std::string& name)
{
	return this->environment->get_root(name);
}

void VM::env_pop()
{
	if (this->env_current_stack != 0) {
		this->env_depth_current_max = this->env_stack_positions[this->env_current_stack].back();
		this->env_stack_positions[this->env_current_stack].pop_back();
	}

	this->env_current_stack = this->env_stack_history.back();
	this->env_stack_history.pop_back();
}

void print_env_stack_pos(std::vector<std::vector<size_t>> env_stack_positions)
{
	auto size = env_stack_positions.size();

	for (size_t i = 0; i < size; i++) {
		std::cout << i << ": [";

		auto size_2 = env_stack_positions[i].size();

		for (size_t ii = 0; ii < size_2; ii++) {
			std::cout << env_stack_positions[i][ii] << ",";
		}

		std::cout << "]\n";
	}
}

void VM::env_push(size_t stack_num)
{
	this->env_stack_history.push_back(this->env_current_stack);
	this->env_current_stack = stack_num;

	if (stack_num == 0) {
		this->environment = this->environment->push();
		return;
	}

	// Initialize vector if neccesary
	if (this->env_current_stack - 2 >= this->env_stack_positions.size()) {
		this->env_stack_positions.push_back(std::vector<size_t>{0});
		this->env_stack_positions.push_back(std::vector<size_t>{0});
		this->env_stack_positions.push_back(std::vector<size_t>{0});
	}
	
	this->env_stack_positions[this->env_current_stack].push_back(this->env_depth_current_max);
}