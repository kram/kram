#include "value.h"

Value::Value() {
	type = Type::NUL;
	number = 0;
	strval = "";
}

Value::Value(Type t) {
	type = t;
	number = 0;
	strval = "";
}

Value::Value(Type t, int val) {
	type = t;
	number = val;
	strval = "";
}

Value::Value(Type t, std::string val) {
	type = t;
	number = 0;
	strval = val;
}

Value* Value::execMethod(std::string name, std::vector<Value*> val) {

	if (this->type != Type::REFERENCE) {
		std::cout << "Is not of type REFERENCE\n";
		exit(0);
	}

	if (this->methods.find(name) == this->methods.end()) {
		std::cout << "UNKNOWN METHOD: " << name << "\n";
		exit(0);
	}

	method m = this->methods[name];

	return m(this, val);
}

void Value::add_method(std::string name, method m) {
	this->methods[name] = m;
}