// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

void VM::gc_add_object(Value* obj) {

	if (obj->type == Type::NUL ||
		obj->type == Type::FUNCTION ||
		obj->type == Type::CLASS) {
		return;
	}

	if (obj->refcount == -1) {
		return;
	}

	// std::cout << "gc_add_object(): " << obj->print() << "\n";

	// Start with refcount at 1
	obj->refcount = 1;

	this->gc_all_objects.push_back(obj);
}

void VM::gc_decrease_refcount(Value* obj) {

	if (obj->type == Type::NUL ||
		obj->type == Type::FUNCTION ||
		obj->type == Type::CLASS) {
		return;
	}

	if (obj->refcount < 0) {
		return;
	}

	--obj->refcount;
}

void VM::gc_decrease_refcount(std::vector<Value*> obj) {
	for (auto &i : obj) {
		this->gc_decrease_refcount(i);
	}
}

void VM::gc_increase_refcount(Value* obj) {

	if (obj->type == Type::NUL ||
		obj->type == Type::FUNCTION ||
		obj->type == Type::CLASS) {
		return;
	}

	if (obj->refcount < 0) {
		return;
	}

	++obj->refcount;
}

void VM::gc_increase_refcount(std::vector<Value*> obj) {
	for (auto &i : obj) {
		this->gc_increase_refcount(i);
	}
}

void VM::gc_clean() {


	// std::cout << "gc_clean()\n";

	std::list<Value*>::iterator it = this->gc_all_objects.begin();

	while ( it != this->gc_all_objects.end() ) {
		//std::cout << (*it)->print() << ": " << (*it)->refcount << "\n";

		if ((*it)->refcount == 0) {

			std::cout << "gc_clean(): deleting " << (*it)->print() << "\n";

			delete *it;
			// *it = nullptr;

			// Erase and get a new iterator
			it = this->gc_all_objects.erase(it);
		} else {
			++it;
		}
	}
}

void VM::gc_print() {
	for (auto &obj : this->gc_all_objects) {
	 	std::cout << obj->print() << ": " << obj->refcount << std::endl;
	}
}