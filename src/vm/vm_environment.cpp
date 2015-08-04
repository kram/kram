// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

void VM::set_name(const char * name, Value* val) {
	this->environment->set(name, val);
}

void VM::set_name_root(const char * name, Value* val) {
	this->environment->set(name, val, true);
}

Value* VM::get_name(const char * name) {
	return this->environment->get(name);
}

Value* VM::get_name_root(const char * name) {
	return this->environment->get(name, true);
}

void VM::env_pop() {
	this->environment = this->environment->pop();
}

void VM::env_push() {
	this->environment = this->environment->push();
}