#include <iostream>
#include "../../value.h"

class IO: public Value {

	static void println(Value* self, Value* val) {
		std::cout << val->print() << "\n";
	}

	public:
		void init() {
			this->add_method("Println", this->println);
		}
};