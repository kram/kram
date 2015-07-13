#include "value.h"

Value::Value() {
	type = Type::NUL;
}

Value::Value(Type t) {
	type = t;
}

Value Value::NUMBER(int val) {
	Value vl(Type::NUMBER);
	vl.number = val;

	return vl;
}

Value Value::STRING(std::string val) {
	Value vl(Type::STRING);
	vl.strval = val;

	return vl;
}