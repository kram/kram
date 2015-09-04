// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include <math.h>

#include "../../value.h"

class Number: public Value {

	static Value* sqrt(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, ::sqrt(val[0]->getNumber()));
	}

	static Value* add(Value* self, std::vector<Value*> val) {

		if (val.size() != 2) {
			std::cout << val.size();
			std::cout << "Number::Add(Number) Excepts exactly 1 parameter\n";
			exit(0);
		}

		if (val[1]->type != Type::NUMBER) {
			std::cout << "Number::Add() Expects the first parameter to be of type Number\n";
			exit(0);
		}

		auto num = val[0]->getNumber();

		num += val[1]->getNumber();

		return new Value(Type::NUMBER, num);

	}

	public:
		void init() {
			this->add_method("Sqrt", this->sqrt);
			this->add_method("Add", this->add);
		}
};