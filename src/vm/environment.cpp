// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "environment.h"
#include <iostream>

Environment::Environment() {
	is_root = false;
	names = new std::unordered_map<std::string, Value*>;
}

void Environment::set(const std::string &name, Value* val) {
	(*this->names)[name] = val;
}

void Environment::set_root(std::string name, Value* val) {
	if (this->is_root) {
		this->set(name, val);
	} else {
		this->root->set(name, val);
	}
}

void Environment::update(const std::string &name, Value* val) {

	if (this->has(name)) {
		(*this->names)[name] = val;
		return;
	}

	if (this->is_root) {
		std::cout << "Unknown name: " << name << "\n";
		exit(0);
	}

	this->parent->update(name, val);
	return;
}

Value* Environment::get(const std::string &name) {

	if (this->has(name)) {
		return (*this->names)[name];
	}

	if (this->is_root) {
		std::cout << "Unknown name: " << name << "\n";
		exit(0);
	}

	return this->parent->get(name);
}

Value* Environment::get_root(const std::string &name) {
	if (this->is_root) {
		return this->get(name);
	}

	return this->root->get(name);
}

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

bool Environment::has(const std::string &name) {

	if (this->names->count(name) == 0) {
		return false;
	}

	return true;
}

bool Environment::exists(const std::string &name) {

	if (this->names->count(name) == 1) {
		return true;
	}

	if (this->is_root) {
		return false;
	}

	return this->parent->exists(name);
}