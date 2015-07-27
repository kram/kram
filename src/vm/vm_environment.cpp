// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

void VM::set_name(const std::string& name, Value* val) {
	this->environment->set(name, val);
}

void VM::set_name_root(std::string name, Value* val) {
	this->environment->set_root(name, val);
}

Value* VM::get_name(const std::string& name) {
	return this->environment->get(name);
}

Value* VM::get_name_root(const std::string& name) {
	return this->environment->get_root(name);
}

void VM::env_pop() {
	this->environment = this->environment->pop();
}

void VM::env_push() {
	this->environment = this->environment->push();
}