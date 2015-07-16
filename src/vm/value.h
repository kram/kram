#ifndef VM_VALUE_H
#define VM_VALUE_H

#include <string>
#include <iostream>
#include <unordered_map>

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

	// Library* ref;
	typedef void (*method)(Value);

	protected:
		std::unordered_map<std::string, method> methods;

	public:
		Type type;

		Value();
		Value(Type);

		static Value NUMBER(int);
		static Value STRING(std::string);
		static Value NUL();
		void REFERENCE(std::string);

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

		/*Library* getReference() {
			return this->ref;
		}*/

		// Overwritten by references
		void init(void) {}

		Value execMethod(std::string name, Value val) {

			if (this->type != Type::REFERENCE) {
				std::cout << "Is not of type REFERENCE\n";
				exit(0);
			}

			std::cout << "Lib::call() " << name << "\n";

			if (this->methods.find(name) == this->methods.end()) {
				std::cout << "UNKNOWN METHOD: " << name << "\n";
				exit(0);
			}

			std::cout << "Pre\n";

			method m = this->methods[name];

			std::cout << "Post\n";

			m(val);

			return Value::NUL();
		}

		void add_method(std::string name, method m) {
			std::cout << "add_method() " << name << "\n";
			this->methods[name] = m;
		}
};

#endif