#include <iostream>
#include "lexer/lexer.h"
#include "parser/parser.h"
#include "vm/instruction.h"
#include "vm/vm.h"

int main() {
	lexer::Lexer lexer;
	std::vector<lexer::Token> tokens = lexer.parse_file();

	// lexer::Lexer::print(tokens);

	Parser parser (tokens);
	std::vector<Instruction> instructions = parser.run();


	std::cout << "Printing ins:\n";

	for (Instruction ins : instructions) {
		ins.print();
		std::cout << "--------\n";
	}

	VM vm;
	vm.boot(instructions);

	return 0;
}