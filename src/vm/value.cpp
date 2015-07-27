// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "value.h"

#include <iostream>

Value::Value() {
	type = Type::NUL;
}

Value::Value(Type t) {
	type = t;

	switch (type) {
		case Type::NUL:
		case Type::BOOL:
		case Type::NUMBER:
			data.number = 0;
			break;

		case Type::STRING:
		case Type::NAME:
			data.strval = new std::string();
			break;

		case Type::REFERENCE:
			data.methods = new Methods();
			break;

		case Type::FUNCTION:
		case Type::CLASS:
			break;
	}
}

Value::Value(Type t, double val) {
	type = t;
	data.number = val;

	switch (type) {
		case Type::NUL:
		case Type::BOOL:
		case Type::NUMBER:
			data.number = val;
			break;

		default:
			std::cout << "Value::Value(Type, double) should not be used with this type!\n";
			exit(0);
			break;
	}
}

Value::Value(Type t, std::string val) {
	type = t;

	switch (type) {
		case Type::STRING:
		case Type::NAME:
			data.strval = new std::string(val);
			break;

		case Type::REFERENCE:
			data.methods = new Methods();
			break;

		case Type::FUNCTION:
		case Type::CLASS:
			break;

		default:
			std::cout << "Value::Value(Type, std::string) should not be used with this type!\n";
			exit(0);
			break;
	}
}

void Value::set_type(Type type) {
	this->type = type;

	switch (type) {
		case Type::NUL:
		case Type::BOOL:
		case Type::NUMBER:
			data.number = 0;
			break;

		case Type::STRING:
		case Type::NAME:
			data.strval = new std::string();
			break;

		case Type::REFERENCE:
			data.methods = new Methods();
			break;

		case Type::FUNCTION:
		case Type::CLASS:
			break;

		default:
			std::cout << "Value::Value(Type) should not be used with this type!\n";
			exit(0);
			break;
	}
}

Value* Value::exec_method(std::string name, std::vector<Value*> val) {

	if (this->type != Type::REFERENCE && this->type != Type::FUNCTION && this->type != Type::CLASS) {
		std::cout << "Can not execute method '" << name << "': Parent is not of type REFERENCE, CLASS, or FUNCTION\n";
		std::cout << this->print() << "\n";	
		exit(0);
	}

	if (this->type == Type::CLASS) {
		std::cout << "Value::exec_method() Should not be used together with Type::CLASS. This is an error, please report it! :)\n";
		exit(0);
	}

	if (this->type == Type::FUNCTION) {
		return this->data.single_method(this, val);
	}

	if (this->data.methods->find(name) == this->data.methods->end()) {
		std::cout << "UNKNOWN METHOD: " << name << "\n";
		exit(0);
	}

	Method m = this->data.methods->at(name);

	return m(this, val);
}

void Value::add_method(std::string name, Method m) {
	this->data.methods->insert( {{name, m}} );
}

bool Value::has_method(std::string name) {

	if (this->data.methods->find(name) == this->data.methods->end()) {
		return false;
	}

	return true;
}