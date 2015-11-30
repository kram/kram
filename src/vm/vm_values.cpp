// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

void VM::init_default_values() {
	this->value_null = new Value(Type::NUL);

	for (int i = 0; i < 255; i++) {
		this->value_number[i] = new Value(Type::NUMBER, i);
	}

	this->value_bool[0] = new Value(Type::BOOL, false);
	this->value_bool[1] = new Value(Type::BOOL, true);
}

Value* VM::get_value_null() {
	return this->value_null;
}

Value* VM::get_value_number(double num) {

	if (num >= 0 && num < 255 && floor(num) == num) {
		return this->value_number[(int) num];
	}

	return new Value(Type::NUMBER, num);
}

Value* VM::get_value_bool(bool val) {
	if (val) {
		return this->value_bool[1];
	}

	return this->value_bool[0];
}
