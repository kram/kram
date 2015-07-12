#include "value.h"

enum class Ins {
	// name, right
	ASSIGN,

	// value
	LITERAL,

	// name
	NAME,

	// left, right, name (the operator)
	MATH,

	// left (true), right (false), center (the if-statement)
	IF,

	// Nothign, indicates an empty instruction in case of failure
	IGNORE,

	// left (the class to push), right (what comes after)
	PUSH_CLASS,

	// left (the method name), right (the parameters)
	CALL
};

class Instruction {
	Ins instruction;
	std::string name;
	Value value;
	std::vector<Instruction> left;
	std::vector<Instruction> right;
	std::vector<Instruction> center;
};