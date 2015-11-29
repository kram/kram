// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include <math.h>

#include "../../value.h"

class Number: public Value {

	static Value* abs(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, std::abs(val[0]->getNumber()));
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

	static Value* pow(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, std::pow(val[0]->getNumber(), val[1]->getNumber()));
	}

	static Value* sqrt(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, std::sqrt(val[0]->getNumber()));
	}

	static Value* sin(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, std::sin(val[0]->getNumber()));
	}

	static Value* cos(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, std::cos(val[0]->getNumber()));
	}

	static Value* tan(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, std::tan(val[0]->getNumber()));
	}

	public:
		void init() {
			this->add_method("Add", this->add);

			this->add_method("Abs", this->abs);
			this->add_method("Pow", this->pow);
			this->add_method("Sqrt", this->sqrt);

			this->add_method("Sin", this->sin);
			this->add_method("Cos", this->cos);
			this->add_method("Tan", this->tan);
		}
};