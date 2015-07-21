#ifndef VM_VALUE_H
#define VM_VALUE_H

#include <string>
#include <iostream>
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
};

class Value {

	typedef Value* (*Method)(Value*, std::vector<Value*>);
	typedef std::unordered_map<std::string, Method> Methods;

	protected:
		union {
			int number;
			std::string* strval;
			Methods* methods;
			Method single_method;
		} data;

	public:

		Type type;

		Value();
		Value(Type);
		Value(Type, std::string);
		Value(Type, int);

		void set_type(Type);

		std::string print(void) {
			std::string res = "";

			std::string i = "UNKNOWN";
			switch (this->type) {
				case Type::NUL: i = "NUL"; break;
				case Type::STRING: i = "STRING"; break;
				case Type::NUMBER: i = "NUMBER"; break;
				case Type::BOOL: i = "BOOL"; break;
				case Type::REFERENCE: i = "REFERENCE"; break;
				case Type::FUNCTION: i = "FUNCTION"; break;
			}

			res += i + "<";

			if (this->type == Type::STRING) {
				res += *this->data.strval;
			}

			if (this->type == Type::NUMBER) {
				res += std::to_string(this->data.number);
			}

			if (this->type == Type::BOOL) {
				if (this->getBool()) {
					res += "true";
				} else {
					res += "false";
				}
			}

			res += ">";

			return res;
		};

		std::string getString() {
			return *this->data.strval;
		}

		int getNumber() {
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
		Value* execMethod(std::string, std::vector<Value*>);
		void add_method(std::string, Method);
};

#endif