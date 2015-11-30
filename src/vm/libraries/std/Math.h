// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "../../value.h"

#include <cmath>

class Math: public Value {

	static Value* pow(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, std::pow(val[0]->getNumber(), val[1]->getNumber()));
	}

	public:

		void init() {
			this->add_method("Pow", this->pow);
		}
};