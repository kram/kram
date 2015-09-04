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
		 * Class::new_instance()
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

		/**
		 * Class::set_value()
		 *
		 * Creates new values, a value with the same name may not exist
		 *
		 * @param std::string name - The key assigned to this value
		 * @param Value* val
		 */
		void set_value(std::string name, Value* val) {
			
			if (this->values.find(name) != this->values.end()) {
				std::cout << "The value, " << name << ", already exists in this Class. Did you mean to use = ?\n";
				exit(0);
			}

			this->values.insert( {{ name, val }} );
		}

		/**
		 * Class::set_value()
		 *
		 * Updates an existing value.
		 * Will be errorous if the value eiter doesn't exist or is of a different type
		 *
		 * @param std::string name - The key assigned to this value
		 * @param Value* val
		 */
		void update_value(std::string name, Value* val) {

			// Not found
			if (this->values.find(name) == this->values.end()) {
				std::cout << "No such value, " << name << ", did you mean to use the := operation?\n";
				exit(0);
			}

			// The values needs to be of the same type
			if (this->values[name]->type != val->type) {
				std::cout << "Can not update (with =) a variable of type " << this->values[name]->print(true) << " to " << val->print(true) << "\n";
				exit(0);
			}

			this->values[name] = val;
		}

		/**
		 * Class::set_value()
		 *
		 * Fetches an existing value. Errors if not set.
		 *
		 * @param std::string name
		 * @return Value*
		 */
		Value* get_value(std::string name) {

			if (this->values.find(name) == this->values.end()) {
				std::cout << "Class has no such value, " << name << "\n";
				exit(0);
			}

			return this->values[name];
		}
};