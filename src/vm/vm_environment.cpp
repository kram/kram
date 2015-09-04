// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

void VM::name_create(const std::string& name, Value* val) {
	this->environment->set(name, val);
}

void VM::name_create_root(std::string name, Value* val) {
	this->environment->set_root(name, val);
}

void VM::name_update(const std::string& name, Value* val) {

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
	return this->environment->get(name);
}

Value* VM::name_get_root(const std::string& name) {
	return this->environment->get_root(name);
}

void VM::env_pop() {
	this->environment = this->environment->pop();
}

void VM::env_push() {
	this->environment = this->environment->push();
}