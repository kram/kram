// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include "../../value.h"

class IO: public Value {

	static Value* debug(Value* self, std::vector<Value*> val) {
		
		for (auto i : val) {
			std::cout << i->print(true) << "\n";
		}

		return new Value();
	}

	static Value* println(Value* self, std::vector<Value*> val) {
		
		for (auto i : val) {
			std::cout << i->print();
		}

		std::cout << "\n";

		return new Value();
	}

	static Value* print(Value* self, std::vector<Value*> val) {
		
		for (auto i : val) {
			std::cout << i->print();
		}

		return new Value();
	}

	public:
		void init() {
			this->add_method("Debug", this->debug);
			this->add_method("Println", this->println);
			this->add_method("Print", this->print);
		}
};