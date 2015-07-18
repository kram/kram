#include <iostream>
#include <math.h>

#include "../../value.h"

class Number: public Value {

	static Value* sqrt(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, ::sqrt(val[0]->getNumber()));
	}

	public:
		void init() {
			this->add_method("Sqrt", this->sqrt);
		}
};