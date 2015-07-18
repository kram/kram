#include <iostream>

#include "lexer/lexer.h"
#include "parser/parser.h"
#include "vm/instruction.h"
#include "vm/vm.h"

int main(int argc, char** argv) {

	bool debug = false;
	std::string file = "";

	// Parse arguments
	// This probable needs to become way more advanced in the future, but it is enough for now...
	if (argc > 1) {
		for (int i = 1; i < argc; i++) {
			std::string arg(argv[i]);

			if (arg == "--debug") {
				debug = true;
			} else {
				file = arg;
			}
		}
	}

	lexer::Lexer lexer;
	std::vector<lexer::Token> tokens = lexer.parse_file(file);

	if (debug) {
		lexer::Lexer::print(tokens);
	}

	Parser parser (tokens);
	std::vector<Instruction> instructions = parser.run();

	if (debug) {
		std::cout << "Printing ins:\n";

		for (Instruction ins : instructions) {
			ins.print();
			std::cout << "--------\n";
		}
	}

	VM vm;
	vm.boot(instructions);

	return 0;
}