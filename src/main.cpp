#include <iostream>
#include "lexer/lexer.h"
#include "parser/parser.h"
#include "vm/instruction.h"

int main() {
	
	std::cout << "Runnign lexer\n";

	lexer::Lexer lexer;
	std::vector<lexer::Token> tokens = lexer.parse_file();

	std::cout << "Lexer result\n";

	for (lexer::Token tok : tokens) {
		tok.print();
	}

	std::cout << "Runnign parser\n";

	Parser parser (tokens);
	std::vector<Instruction> instructions = parser.run();

	std::cout << "Parser result\n";

	for (int i = 0; i < instructions.size() - 1; i++) {
		instructions[i].print();
	}

	return 0;
}