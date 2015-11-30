// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "../../output.h"
#include "../../value.h"

extern Output * kram_output_stream;

class IO: public Value {

	static Value* debug(Value* self, std::vector<Value*> val) {
		
		for (auto i : val) {
			(*kram_output_stream) << i->print(true) << "\n";
		}

		return new Value();
	}

	static Value* println(Value* self, std::vector<Value*> val) {
		
		for (auto i : val) {
			(*kram_output_stream) << i->print();
		}

		(*kram_output_stream) << "\n";

		return new Value();
	}

	static Value* print(Value* self, std::vector<Value*> val) {
		
		for (auto i : val) {
			(*kram_output_stream) << i->print();
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