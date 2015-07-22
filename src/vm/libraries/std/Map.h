#include <iostream>
#include <unordered_map>

#include "../../value.h"

class Map: public Value {

	static Value* set(Value* self, std::vector<Value*> val) {
		Map* map = static_cast<Map*>(self);

		if (val.size() != 2) {
			std::cout << "Map.Set() expects exacly two parameters\n";
			exit(0);
		}

		if (val[0]->type != Type::STRING) {
			
		}

		map.insert( {{  }} );

		return new Value();
	}

	static Value* get(Value* self, std::vector<Value*> val) {
		return new Value();
	}

	static Value* has(Value* self, std::vector<Value*> val) {
		return new Value();
	}

	public:
		unordered_map<std::string, Value*> content;

		void init() {
			this->add_method("Set", this->add);
			this->add_method("Get", this->get);
			this->add_method("Has", this->has);
		}
};