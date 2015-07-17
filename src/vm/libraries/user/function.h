#include <iostream>
#include "../../value.h"
#include "../../vm.h"

class Function: public Value {

	static void exec(Value* self, Value* val) {
		Function* fn = static_cast<Function*>(self);
		fn->vm->run(fn->content);
	}

	public:
		std::vector<Instruction> content;
		VM* vm;

		void init() {
			this->add_method("exec", this->exec);
		}

		void set_content(std::vector<Instruction> ins) {
			this->content = ins;
		}
};