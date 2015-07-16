#include <iostream>
#include "../../value.h"
#include "../../vm.h"

class Function: public Value {

	static void exec(Value val) {
		std::cout << "Executing function\n";
		//Function* fn = (Function*) self;
		//Function* fn = static_cast<Function*>(self);
		// self->vm->run(fn->content);
	}

	public:
		std::vector<Instruction> content;

		void init() {
			this->add_method("exec", this->exec);
		}

		void set_content(std::vector<Instruction> ins) {
			this->content = ins;
		}
};