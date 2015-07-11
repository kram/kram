#include <iostream>
#include "lexer/lexer.h"

int main() {
	
	Lexer lex;
	std::vector<Token> tokens = lex.parse_file();

	for (Token tok : tokens) {
		tok.print();
	}

	return 0;
}