// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include "../../value.h"

class IO: public Value {

	static Value* println(Value* self, std::vector<Value*> val) {
		std::cout << val[0]->print() << "\n";
		return new Value();
	}

	public:
		void init() {
			this->add_method("Println", this->println);
		}
};