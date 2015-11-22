// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include "../../value.h"
#include "../../vm.h"

class Function: public Value {

	static Value* exec(Value* self, std::vector<Value*> arguments) {
		Function* fn = static_cast<Function*>(self);

		if (arguments.size() < fn->required_parameter_count) {
			std::cout << "Got fewer arguments than the required count\n";
			exit(1);
		}

		// Assign function parameters as variables
		size_t key = 0;
		for (Instruction* param : fn->parameters) {

			if (arguments.size() > key) {

				// Caller defined value
				fn->vm->name_create(param->name, arguments[key]);

			} else if (param->right.size() == 1) {

				// Use default value
				fn->vm->name_create(param->name, fn->vm->run(param->right[0]));

			} else {
				std::cout << "No value for parameter '" << param->name << "'\n";
				exit(1);
			}

			++key;
		}

		return fn->vm->run(fn->content);
	}

	public:
		std::vector<Instruction*> content;
		std::vector<Instruction*> parameters;
		size_t required_parameter_count;
		VM* vm;

		void init() {
			this->data.single_method = this->exec;
		}

		void set_content(std::vector<Instruction*> ins) {
			this->content = ins;
		}

		void set_parameters(std::vector<Instruction*> ins) {
			this->parameters = ins;

			this->required_parameter_count = 0;

			// Pre-calculated values and validation
			for (Instruction* i : ins) {

				if (i->instruction != Ins::FUNCTION_PARAMETER) {
					std::cout << "A function can only take parameters of type FUNCTION_PARAMETER\n";
					exit(1);
				}

				if (i->right.size() == 0) {
					this->required_parameter_count++;
					continue;
				}

				if (i->right.size() > 1) {
					std::cout << "Function: The size of a default parameter can only be one!\n";
					exit(1);
				}
			}
		}
};