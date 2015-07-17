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