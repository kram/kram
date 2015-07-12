#include <iostream>
#include <fstream>
#include <string>

#include "lexer.h"

Lexer::Lexer() {
	// Initialize all operators and keywords
	operators["+"] = true;
	operators["-"] = true;
	operators["*"] = true;
	operators["/"] = true;
	operators["%"] = true;
	operators["**"] = true;
	operators["="] = true;
	operators["=="] = true;
	operators[">"] = true;
	operators[">="] = true;
	operators["<"] = true;
	operators["<="] = true;
	operators["&&"] = true;
	operators["|"] = true;
	operators["||"] = true;
	operators["..."] = true;
	operators[".."] = true;
	operators["."] = true;
	operators["{"] = true;
	operators["}"] = true;
	operators[":"] = true;
	operators["++"] = true;
	operators["--"] = true;

	keywords["if"] = true;
	keywords["else"] = true;
	keywords["var"] = true;
	keywords["class"] = true;
	keywords["static"] = true;
	keywords["return"] = true;
	keywords["for"] = true;
	keywords["in"] = true;
}

std::vector<Token> Lexer::parse_file() {
	std::vector<Token> result;

	std::ifstream file("test.kr");

    while (std::getline(file, this->row))
    {
    	this->index = 0;

    	while (true) {
        	Token tok = this->next();

        	this->index++;

        	if (tok.type == Type::IGNORE) {
        		continue;
        	}

        	result.push_back(tok);

        	if (tok.type == Type::T_EOF || tok.type == Type::T_EOL) {
        		break;
        	}
        }
    }

    // Indicate end of file
    result.push_back(Token::T_EOF());

    return result;
}

char Lexer::char_at_pos(int index) {
	// Indicate nothingness
	if (index >= this->row.size()) {
		return '\0';
	}

	return this->row[index];
}

Token Lexer::next() {
	this->current = this->char_at_pos(this->index);

	// End of row
	if (this->current == '\0') {
		return Token::T_EOL();
	}

	// Ignore Whitespace
	if (this->current == ' ') {
		return Token::IGNORE();
	}

	// Comments
	if (this->current == '/' && this->char_at_pos(this->index+1) == '/') {
		return this->comment();
	}

	// Names
	// Begins with a char a-Z
	if ((this->current >= 'a' && this->current <= 'z') || (this->current >= 'A' && this->current <= 'Z')) {
		return this->name();
	}

	// Numbers
	if (this->current >= '0' && this->current <= '9') {
		return this->number();
	}

	// Strings
	if (this->current == '\"') {
		return this->string();
	}

	// operators
	if (this->operators.find(&this->current) != this->operators.end()) {
		return this->oper();
	}

	return Token::IGNORE();
}

Token Lexer::comment() {
	while(true) {
		this->index += 1;
		char current = this->char_at_pos(this->index + 1);

		if (current == '\n' || current == '\r') {
			break;
		}
	}

	return Token::T_EOL();
}

Token Lexer::name() {
	std::string s = &this->current;

	while (true) {
		char c = this->char_at_pos(this->index + 1);

		// After the beginning, a name can be a-Z0-9_
		if ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			s += c;
			this->index += 1;
		} else {
			break;
		}
	}

	if (s == "true" || s == "false") {
		return Token::BOOL(s);
	}

	if (this->keywords.find(s) != this->keywords.end()) {
		return Token::KEYWORD(s);
	}

	return Token::NAME(s);
}

Token Lexer::number() {
	std::string s = &this->current;

	// Look for more digits.
	while(true) {
		char c = this->char_at_pos(this->index + 1);

		if ((c < '0' || c > '9') && c != '.') {
			break;
		}

		// A dot needs to be followed by another digit to be valid
		if (c == '.') {
			char cc = this->char_at_pos(this->index + 2);

			if (cc < '0' || cc > '9') {
				break;
			}
		}

		this->index += 1;
		s += c;
	}

	// TODO Decimal
	// TODO Verify that it ends with a space?

	return Token::NUMBER(s);
}

Token Lexer::string() {
	std::string s = &this->current;

	this->index += 1;

	while(true) {

		// End of string
		if (this->char_at_pos(this->index) == '"') {
			break;
		}

		// Escaping
		if (this->char_at_pos(this->index) == '\\') {
			this->index += 1;
		}

		s += this->char_at_pos(this->index);
		this->index += 1;
	}

	return Token::STRING(s);
}

Token Lexer::oper() {
	std::string s = &this->current;

	while(true) {

		char next = this->char_at_pos(this->index + 1);

		// EOF
		if (next == '\0') {
			break;
		}

		std::string combined = s + next;

		if (this->keywords.find(combined) != this->keywords.end()) {
			s += next;
			this->index += 1;
		} else {
			break;
		}
	}

	return Token::OPERATOR(s);
}
