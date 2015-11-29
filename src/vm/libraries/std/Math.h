// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "../../value.h"

class Math: public Value {

	static Value* kr_abs(Value* self, std::vector<Value*> val) {

		auto num = val[0]->getNumber();

		return new Value(Type::NUMBER, std::abs(num));
	}

	public:

		void init() {
			this->add_method("Abs", this->kr_abs);
		}
};