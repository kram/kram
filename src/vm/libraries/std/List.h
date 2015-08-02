// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include <vector>

#include "../../value.h"

class List: public Value {

	static Value* constructor(Value* self, std::vector<Value*> val) {
		List* instance = new List();
		instance->set_type(Type::REFERENCE);
		instance->init();

		return instance;
	}

	static Value* push(Value* self, std::vector<Value*> val) {

		if (val.size() == 0) {
			std::cout << "List.Push() expects at least 1 parameter\n";
			exit(0);
		}

		List* list = static_cast<List*>(self);

		for (Value* i : val) {
			list->content.push_back(i);
		}

		return list;
	}

	static Value* pop(Value* self, std::vector<Value*> val) {

		if (val.size() != 0) {
			std::cout << "List.Pop() takes no parameters\n";
			exit(0);
		}

		List* list = static_cast<List*>(self);

		if (list->content.size() == 0) {
			std::cout << "List.Pop() requires the list to contain at least 1 item\n";
			exit(0);
		}

		Value* res = list->content.back();

		list->content.pop_back();

		return res;
	}

	static Value* at(Value* self, std::vector<Value*> val) {

		if (val.size() == 0) {
			std::cout << "List.At() expects exactly 1 parameter\n";
			exit(0);
		}

		if (val[0]->type != Type::NUMBER) {
			std::cout << "List.At() expects parameter 1 to be of type NUMBER\n";
			exit(0);
		}

		List* list = static_cast<List*>(self);

		size_t at = val[0]->getNumber();

		if (at < 0 || at >= list->content.size()) {
			std::cout << "List.At() out of bounds\n";
			exit(0);
		}

		return list->content.at(at);
	}

	static Value* size(Value* self, std::vector<Value*> val) {

		if (val.size() != 0) {
			std::cout << "List.Size() takes no parameters\n";
			exit(0);
		}

		List* list = static_cast<List*>(self);
		
		return new Value(Type::NUMBER, list->content.size());
	}

	public:
		std::vector<Value*> content;

		void init() {
			this->add_method("new", this->constructor);

			this->add_method("Push", this->push);
			this->add_method("Pop", this->pop);
			this->add_method("At", this->at);
			this->add_method("Size", this->size);
		}
};