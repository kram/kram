#include <iostream>
#include "../../value.h"
#include "../../vm.h"

class Function: public Value {

	static Value* exec(Value* self, std::vector<Value*> val) {
		Function* fn = static_cast<Function*>(self);

		if (val.size() != fn->parameters.size()) {
			std::cout << "Argument and parameter count needs to match\n";
			exit(0);
		}

		int key = 0;
		for (Instruction* param : fn->parameters) {
			fn->vm->set_name(param->name, val[key]);
			++key;
		}

		return fn->vm->run(fn->content);
	}

	public:
		std::vector<Instruction*> content;
		std::vector<Instruction*> parameters;
		VM* vm;

		void init() {
			this->data.single_method = this->exec;
		}

		void set_content(std::vector<Instruction*> ins) {
			this->content = ins;
		}

		void set_parameters(std::vector<Instruction*> ins) {
			this->parameters = ins;
		}
};