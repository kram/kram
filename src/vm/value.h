#ifndef VM_VALUE_H
#define VM_VALUE_H

#include <string>
#include <iostream>

// Fake library
class Library;

enum class Type {
	NUL,
	STRING,
	NUMBER,
	REFERENCE
};

class Value {
	Type type;
	
	// TODO: Union(-ify) this
	std::string strval;
	int number;
	Library* ref;

	public:
		Value();
		Value(Type);

		static Value NUMBER(int);
		static Value STRING(std::string);
		static Value REFERENCE(Library*);
		static Value NUL();

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

			res += ">";

			return res;
		};

		std::string string() {
			return this->strval;
		}
};

#endif