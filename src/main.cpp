#include <iostream>
#include "lexer/lexer.h"
#include "parser/parser.h"

int main() {
	
	lexer::Lexer lexer;
	std::vector<lexer::Token> tokens = lexer.parse_file();

	for (lexer::Token tok : tokens) {
		tok.print();
	}

	Parser parser (tokens);

	return 0;
}