#include <iostream>
#include <fstream>
#include <string>

#include "lexer.h"

using namespace lexer;

Lexer::Lexer() {
	// Initialize all  keywords
	keywords["class"] = true;
	keywords["fn"] = true;
	keywords["if"] = true;
	keywords["else"] = true;
	keywords["new"] = true;
}

std::vector<Token*> Lexer::parse_file(std::string filename) {
	std::vector<Token*> result;

	std::ifstream file(filename);

	if (file.fail()) {
		std::cout << "Could not open " << filename << "\n";
		exit(0);
	}

    while (std::getline(file, this->row))
    {
    	this->index = 0;

    	while (true) {
        	Token* tok = this->next();

        	this->index++;

        	if (tok->type == Type::IGNORE) {
        		continue;
        	}

        	result.push_back(tok);

        	if (tok->type == Type::T_EOF || tok->type == Type::T_EOL) {
        		break;
        	}
        }
    }

    // Indicate end of file
    result.push_back(new Token(Type::T_EOF));

    return result;
}

char Lexer::char_at_pos(size_t index) {
	// Indicate nothingness
	if (index >= this->row.size()) {
		return '\0';
	}

	return this->row[index];
}

Token* Lexer::next() {
	this->current = this->char_at_pos(this->index);

	// End of row
	if (this->current == '\0') {
		return new Token(Type::T_EOL);
	}

	// Ignore Whitespace
	if (iswspace(this->current)) {
		return new Token(Type::IGNORE);
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

	// Operators
	std::string str(1, this->current);
	if (Token::Trans(str) != Type::IGNORE) {
		return this->oper();
	}

	std::cout << "Ignoring: " << this->current << "(" << str << ")\n";

	return new Token(Type::IGNORE);
}

Token* Lexer::comment() {
	while(true) {
		this->index += 1;
		char current = this->char_at_pos(this->index + 1);

		if (current == '\n' || current == '\r' || current == '\0') {
			break;
		}
	}

	return new Token(Type::T_EOL);
}

Token* Lexer::name() {
	std::string s(1, this->current);

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
		return new Token(Type::BOOL, s);
	}

	if (this->keywords.find(s) != this->keywords.end()) {
		return new Token(Type::KEYWORD, s);
	}

	return new Token(Type::NAME, s);
}

Token* Lexer::number() {
	std::string s(1, this->current);

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

	return new Token(Type::NUMBER, s);
}

Token* Lexer::string() {
	std::string s = "";

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

	return new Token(Type::STRING, s);
}

Token* Lexer::oper() {
	std::string s(1, this->current);

	while(true) {

		char next = this->char_at_pos(this->index + 1);

		// EOF
		if (next == '\0') {
			break;
		}

		std::string combined = s + next;

		if (Token::Trans(combined) != Type::IGNORE) {
			s += next;
			this->index += 1;
		} else {
			break;
		}
	}

	return new Token(Type::OPERATOR, s);
}

void Lexer::print(std::vector<Token*> tokens) {
	for (Token* tok : tokens) {
		std::cout << tok->print() << "\n";
	}
}