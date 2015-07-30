// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include "../../value.h"
#include "../../vm.h"

class Class: public Value {

	public:
		VM* vm;

		std::unordered_map<std::string, Value*> values;

		/**
		 * Value::new_instance()
		 *
		 * Creates a new Class and clones this->values
		 */
		Value* new_instance() {

			Class* cl = new Class();
			cl->type = Type::CLASS;
			cl->vm = this->vm;

			for (auto i : this->values) {
				cl->values.insert( {{ i.first, i.second->clone() }} );
			}

			return cl;
		}

		void set_value(std::string name, Value* val) {
			// Replace or create new
			if (this->values.find(name) == this->values.end()) {
				this->values.insert( {{ name, val }} );
			} else {
				this->values[name] = val;
			}
		}

		Value* get_value(std::string name) {

			if (this->values.find(name) == this->values.end()) {
				std::cout << "Class has no such value, " << name << "\n";
				exit(0);
			}

			return this->values[name];
		}
};