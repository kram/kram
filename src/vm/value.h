// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#ifndef VM_VALUE_H
#define VM_VALUE_H

#include <string>
#include <sstream>
#include <unordered_map>
#include <vector>

class VM;

enum class Type {
	NUL,
	STRING,
	NUMBER,
	BOOL,
	REFERENCE,
	FUNCTION,
	CLASS,
	NAME,
};

class Value {

	typedef Value* (*Method)(Value*, std::vector<Value*>);
	typedef std::unordered_map<std::string, Method> Methods;

	protected:
		union {
			double number;
			std::string* strval;
			Methods* methods;
			Method single_method;
		} data;

	public:

		Type type;

		Value();
		Value(Type);
		Value(Type, std::string);
		Value(Type, double);

		void set_type(Type);

		std::string print(void) {
			std::stringstream res;

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

			if (this->type == Type::STRING) {
				res << *this->data.strval;
			}

			if (this->type == Type::NUMBER) {
				res << this->data.number;
			}

			if (this->type == Type::BOOL) {
				if (this->getBool()) {
					res << "true";
				} else {
					res << "false";
				}
			}

			/*if (this->type == Type::CLASS) {
				Class* cl = static_cast<Class*>(this);
				res << cl->print_values();
			}*/

			res << ">";

			return res.str();
		};

		std::string getString() {
			return *this->data.strval;
		}

		double getNumber() {
			return this->data.number;
		}

		bool getBool() {
			if (this->data.number == 0) {
				return false;
			}

			return true;
		}

		// Overwritten by references
		void init(void) {}

		// #justlibrarythings
		Value* exec_method(std::string, std::vector<Value*>);
		void add_method(std::string, Method);
		bool has_method(std::string);
};

#endif