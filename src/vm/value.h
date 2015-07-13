#include <string>
#include <iostream>

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

	public:
		Value();
		Value(Type);
		static Value NUMBER(int);
		static Value STRING(std::string);

		void print(void) {
			std::string i = "UNKNOWN";

			switch (this->type) {
				case Type::NUL: i = "NUL"; break;
				case Type::STRING: i = "STRING"; break;
				case Type::NUMBER: i = "NUMBER"; break;
				case Type::REFERENCE: i = "REFERENCE"; break;
			}

			std::cout << i << "<";

			if (this->type == Type::STRING) {
				std::cout << this->strval;
			}

			if (this->type == Type::NUMBER) {
				std::cout << this->number;
			}

			std::cout << ">\n";
		};
};