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

std::string Value::print(bool print_type) {
	std::stringstream res;

	if (print_type) {
		std::string i = "UNKNOWN";
		switch (this->type) {
			case Type::NUL: i = "NUL"; break;
			case Type::STRING: i = "STRING"; break;
			case Type::NUMBER: i = "NUMBER"; break;
			case Type::BOOL: i = "BOOL"; break;
			case Type::REFERENCE: i = "REFERENCE"; break;
			case Type::FUNCTION: i = "FUNCTION"; break;
			case Type::CLASS: i = "CLASS"; break;
			case Type::NAME: i = "NAME"; break;
		}

		res << i << "<";
	}

	switch (this->type) {
		case Type::NUL:
			res << "NULL";
			break;

		case Type::STRING:
			res << *this->data.strval;
			break;

		case Type::NUMBER: 
			res << this->data.number;
			break;

		case Type::BOOL:
			if (this->getBool()) {
				res << "true";
			} else {
				res << "false";
			}
			break;

		// We can probably prettify these even further, eg print what a REFERENCE is actually pointing to.
		// But that is for another time.
		case Type::REFERENCE: res << "REFERENCE"; break;
		case Type::FUNCTION: res << "FUNCTION"; break;
		case Type::CLASS: res << "CLASS"; break;
		case Type::NAME: res << "NAME"; break;
	}

	if (print_type) {
		res << ">";
	}

	return res.str();
};

std::string Value::getString() {
	return *this->data.strval;
}

double Value::getNumber() {
	return this->data.number;
}

bool Value::getBool() {
	if (this->data.number == 0) {
		return false;
	}

	return true;
}

Value* Value::clone() {

	Value* val = new Value();
	val->type = this->type;

	switch (val->type) {
		case Type::NUMBER:
		case Type::BOOL:
			val->data.number = this->data.number;
			break;

		case Type::STRING:
			val->data.strval = new std::string(*this->data.strval);
			break;

		case Type::REFERENCE:
		case Type::FUNCTION:
		case Type::CLASS:
			delete val;
			return this;
			break;

		default:
			std::cout << "Unable to clone value\n";
			exit(0);
			break;
	}

	return val;
}