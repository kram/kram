#include "environment.h"
#include <iostream>

Environment::Environment() {
	is_root = false;
}

void Environment::set(std::string name, Value* val) {
	this->names[name] = val;
	this->all_names[name] = val;
}

void Environment::set_root(std::string name, Value* val) {
	if (this->is_root) {
		this->set(name, val);
	} else {
		this->root->set(name, val);
	}
}

Value* Environment::get(std::string name) {

	if (!this->has(name)) {
		std::cout << "Unknown name: " << name << "\n";
		exit(0);
	}

	return this->all_names[name];
}

Value* Environment::get_root(std::string name) {
	if (this->is_root) {
		return this->get(name);
	}

	return this->root->get(name);
}

Environment* Environment::push() {
	Environment* env = new Environment();
	env->parent = this;
	env->all_names = this->all_names;

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

bool Environment::has(std::string name) {

	if (this->all_names.find(name) == this->names.end()) {
		return false;
	}

	return true;
}