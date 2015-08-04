// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "environment.h"
#include <iostream>

Environment::Environment() {
	is_root = false;
	names = new Kram_Map::map;
}

void Environment::set(const char * name, Value* val, bool root) {

	if (root && !this->is_root) {
		return this->root->set(name, val);
	}

	(*this->names)[name] = val;
}

/*void Environment::set_root(const char * name, Value* val) {
	if (this->is_root) {
		this->set(name, val);
	} else {
		this->root->set(name, val);
	}
}
*/
Value* Environment::get(const char * name, bool root) {

	if (root && !this->is_root) {
		return this->root->get(name);
	}

	if (this->has(name)) {
		return (*this->names)[name];
	}

	if (this->is_root) {
		std::cout << "Unknown name: " << name << "\n";
		exit(0);
	}

	return this->parent->get(name);
}

/*Value* Environment::get_root(const char * name) {
	if (this->is_root) {
		return this->get(name);
	}

	return this->root->get(name);
}*/

Environment* Environment::push() {
	Environment* env = new Environment();
	env->parent = this;

	if (this->is_root) {
		env->root = this;
	} else {
		env->root = this->root;
	}

	return env;
}

Environment* Environment::pop() {
	if (this->is_root) {
		std::cout << "Environment: You can not pop the root level!" << "\n";
		exit(0);
	}

	return this->parent;
}

bool Environment::has(const char * name) {

	if (this->names->count(name) == 0) {
		return false;
	}

	return true;
}