#ifndef VM_VALUE_H
#define VM_VALUE_H

#include <string>
#include <iostream>
#include <unordered_map>

class VM;

enum class Type {
	NUL,
	STRING,
	NUMBER,
	REFERENCE
};

class Value {
	
	// TODO: Union(-ify) this
	std::string strval;
	int number;

	typedef void (*method)(Value*, Value*);

	protected:
		std::unordered_map<std::string, method> methods;

	public:
		Type type;

		Value();
		Value(Type);
		Value(Type, std::string);
		Value(Type, int);

		std::string print(void) {
			std::string res = "";

			std::string i = "UNKNOWN";
			switch (this->type) {
				case Type::NUL: i = "NUL"; break;
				case Type::STRING: i = "STRING"; break;
				case Type::NUMBER: i = "NUMBER"; break;
				case Type::REFERENCE: i = "REFERENCE"; break;
			}

			res += i + "<";

			if (this->type == Type::STRING) {
				res += this->strval;
			}

			if (this->type == Type::NUMBER) {
				res += std::to_string(this->number);
			}

			if (this->type == Type::REFERENCE) {
				res += this->strval;
				//res += this->ref->print();
			}

			res += ">";

			return res;
		};

		std::string getString() {
			return this->strval;
		}

		int getNumber() {
			return this->number;
		}

		// Overwritten by references
		void init(void) {}

		// #justlibrarythings
		Value* execMethod(std::string name, Value* val);
		void add_method(std::string name, method m);
};

#endif