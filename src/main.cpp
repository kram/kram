#include <iostream>
#include "lexer/lexer.h"
#include "parser/parser.h"
#include "vm/instruction.h"
#include "vm/vm.h"

int main() {
	lexer::Lexer lexer;
	std::vector<lexer::Token> tokens = lexer.parse_file();

	Parser parser (tokens);
	std::vector<Instruction> instructions = parser.run();

	for (int i = 0; i < instructions.size(); i++) {
		instructions[i].print();
	}

	VM vm;
	vm.boot(instructions);

	return 0;
}