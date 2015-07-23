#include <iostream>
#include <unordered_map>

#include "../../value.h"

class Map: public Value {

	static Value* constructor(Value* self, std::vector<Value*> val) {
		Map* instance = new Map();
		instance->set_type(Type::REFERENCE);
		instance->init();

		return instance;
	}

	static Value* set(Value* self, std::vector<Value*> val) {
		Map* map = static_cast<Map*>(self);

		if (val.size() != 2) {
			std::cout << "Map.Set() expects exacly two parameters\n";
			exit(0);
		}

		if (val[0]->type != Type::STRING) {
			std::cout << "Map.Set() expects first parameter to be of type String\n";
			exit(0);
		}

		map->content.insert( {{ val[0]->getString(), val[1] }} );

		return val[1];
	}

	static Value* get(Value* self, std::vector<Value*> val) {
		Map* map = static_cast<Map*>(self);

		if (map->content.find(val[0]->getString()) == map->content.end()) {
			std::cout << "Map.Get() unknown key, \"" << val[0]->getString() << "\"\n";
			exit(0);
		}

		return map->content.at(val[0]->getString());
	}

	static Value* has(Value* self, std::vector<Value*> val) {
		Map* map = static_cast<Map*>(self);

		if (map->content.find(val[0]->getString()) == map->content.end()) {
			return new Value(Type::BOOL, 0);
		}

		return new Value(Type::BOOL, 1);
	}

	public:
		std::unordered_map<std::string, Value*> content;

		void init() {
			this->add_method("new", this->constructor);

			this->add_method("Set", this->set);
			this->add_method("Get", this->get);
			this->add_method("Has", this->has);
		}
};