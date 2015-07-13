#ifndef VM_INSTRUCTION_H
#define VM_INSTRUCTION_H

#include <string>
#include <iostream>
#include <vector>
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

	// Nothing, indicates an empty instruction in case of failure
	IGNORE,

	// left (the class to push), right (what comes after)
	PUSH_CLASS,

	// left (the method name), right (the parameters)
	CALL
};

class Instruction {
	public:
		Ins instruction;
		std::string name;
		Value value;
		std::vector<Instruction> left;
		std::vector<Instruction> right;
		std::vector<Instruction> center;

		Instruction(Ins i) : instruction(i) {}

		void print(int ident = 0) {
			std::string i = "UNKNOWN";

			switch (this->instruction) {
				case Ins::ASSIGN: i = "ASSIGN"; break;
				case Ins::LITERAL: i = "LITERAL"; break;
				case Ins::NAME: i = "NAME"; break;
				case Ins::MATH: i = "MATH"; break;
				case Ins::IF: i = "IF"; break;
				case Ins::IGNORE: i = "IGNORE"; break;
				case Ins::PUSH_CLASS: i = "PUSH_CLASS"; break;
				case Ins::CALL: i = "CALL"; break;
			}

			std::cout << std::string(ident, '\t') << "{\n";
			std::cout << std::string(ident + 1, '\t') << "instruction: " << i << "\n";
			std::cout << std::string(ident + 1, '\t') << "name: " << this->name << "\n";
			std::cout << std::string(ident + 1, '\t') << "value: "; this->value.print();

			if (left.size() > 0) {
				std::cout << std::string(ident + 1, '\t') << "left: [\n";

				for (Instruction i : left) {
					i.print(ident+1);
				}

				std::cout << std::string(ident + 1, '\t') << "]\n";
			}

			/*if (center.size() > 0) {
				std::cout << std::string(ident + 1, '\t') << "center: [\n";

				for (Instruction i : center) {
					i.print(ident+1);
				}

				std::cout << std::string(ident + 1, '\t') << "]\n";
			}*/

			if (right.size() > 0) {
				std::cout << std::string(ident + 1, '\t') << "right: [\n";

				for (Instruction i : right) {
					i.print(ident+1);
				}

				std::cout << std::string(ident + 1, '\t') << "]\n";
			}

			std::cout << std::string(ident, '\t') << "}\n";
		}
};

#endif