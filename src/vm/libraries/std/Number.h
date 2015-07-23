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

	public:
		void init() {
			this->add_method("Sqrt", this->sqrt);
		}
};