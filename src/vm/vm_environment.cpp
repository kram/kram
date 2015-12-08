// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

#include <iostream>

size_t VM::env_get_pos(stack_and_pos snp) {
	size_t stack = KR_SNP_GET_STACK(snp);
	size_t pos = KR_SNP_GET_POS(snp);

	size_t stack_start = this->env_stack_positions[stack].back();

	return stack_start + pos;
}

void VM::name_create(const std::string& name, Value* val) {

	// std::cout << "name_create: " << name << "\n";

	this->environment->set(name, val);
}

void VM::name_create(size_t pos, Value* val) {
	// std::cout << "name_create: " << pos << "\n";

	if (pos > this->env_depth_current_max) {
		this->env_depth_current_max = pos;
	}

	this->environment->set(pos, val);
}

void VM::name_create_root(std::string name, Value* val) {
	this->environment->set_root(name, val);
}

void VM::name_update(size_t pos, Value* val) {
	// std::cout << "name_update: " << pos << "\n";
	this->environment->update(pos, val);
}

void VM::name_update(const std::string& name, Value* val) {
	std::cout << "name_update: " << name << "\n";

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

Value* VM::name_get(const std::string& name) {
	// std::cout << "name_get: " << name << "\n";
	return this->environment->get(name);
}

Value* VM::name_get(size_t pos) {
	//std::cout << "name_get: " << pos << "\n";
	return this->environment->get(pos);
}

Value* VM::name_get_root(const std::string& name) {
	return this->environment->get_root(name);
}

void VM::env_pop() {
	//std::cout << "env_pop" << "\n";

	if (this->env_current_stack != 0) {
		this->env_depth_current_max = this->env_stack_positions[this->env_current_stack].back();
		this->env_stack_positions[this->env_current_stack].pop_back();
	}

	this->env_current_stack = this->env_stack_history.back();
	this->env_stack_history.pop_back();

	//this->env_depth_pos = this->env_depth_pos_stack.back();
	// this->env_depth_pos_stack.pop_back();

	//this->environment = this->environment->pop();
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

void VM::env_push(size_t stack_num) {

	// std::cout << "env_push (" << stack_num << ")\n";

	// std::cout << "this->env_current_stack = " << this->env_current_stack << "\n"; 
	// std::cout << "this->env_stack_positions.size() = " << this->env_stack_positions.size() << "\n";

	// std::cout << "print 1\n";
	// print_env_stack_pos(this->env_stack_positions);
	// std::cout << "print 1 (done)\n";

	this->env_stack_history.push_back(this->env_current_stack);
	this->env_current_stack = stack_num;

	if (stack_num == 0) {
		this->environment = this->environment->push();
		return;
	}

	// this->env_current_stack++;

	// Initialize vector if neccesary
	if (this->env_current_stack >= this->env_stack_positions.size()) {
		this->env_stack_positions.push_back(std::vector<size_t>{0});
	}

	//std::cout << "print 2\n";
	//print_env_stack_pos(this->env_stack_positions);
	//std::cout << "print 2 (done)\n";
	
	this->env_stack_positions[this->env_current_stack].push_back(this->env_depth_current_max);

	//std::cout << "print\n";
	//print_env_stack_pos(this->env_stack_positions);
	//std::cout << "print (done)\n";

	//this->env_depth_pos_stack.push_back(this->env_depth_current_max);
	//this->env_depth_pos = this->env_depth_current_max;
	//this->environment = this->environment->push();

	//std::cout << "env_push (done)" << "\n";
	// exit(1);
}