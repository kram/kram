#include <iostream>
#include "token.h"

using namespace lexer;

Token::Token() {
	type = Type::IGNORE;
	sub = Type::IGNORE;
}

Token Token::T_EOF() {
	Token tok;
	tok.type = Type::T_EOF;
	return tok;
}

Token Token::T_EOL() {
	Token tok;
	tok.type = Type::T_EOL;
	return tok;
}

Token Token::IGNORE() {
	Token tok;
	tok.type = Type::IGNORE;
	return tok;
}

Token Token::STRING(std::string val) {
	Token tok;
	tok.type = Type::STRING;
	tok.value = val;
	return tok;
}

Token Token::NUMBER(std::string val) {
	Token tok;
	tok.type = Type::NUMBER;
	tok.value = val;
	return tok;
}

Token Token::KEYWORD(std::string val) {
	Token tok;
	tok.type = Type::KEYWORD;
	tok.value = val;

	if (val == "var") {
		tok.sub = Type::KEYWORD_VAR;
	}

	return tok;
}

Token Token::OPERATOR(std::string val) {
	Token tok;
	tok.type = Type::OPERATOR;
	tok.value = val;

	if (val == "=") {
		tok.sub = Type::OPERATOR_EQ;
	} else if (val == ".") {
		tok.sub = Type::OPERATOR_DOT;
	} else if (val == "(") {
		tok.sub = Type::OPERATOR_PAREN_L;
	} else if (val == ")") {
		tok.sub = Type::OPERATOR_PAREN_R;
	} else if (val == ";") {
		tok.sub = Type::OPERATOR_SEMICOLON;
	} else if (val == ",") {
		tok.sub = Type::OPERATOR_COLON;
	}

	return tok;
}

Token Token::NAME(std::string val) {
	Token tok;
	tok.type = Type::NAME;
	tok.value = val;
	return tok;
}

Token Token::BOOL(std::string val) {
	Token tok;
	tok.type = Type::BOOL;
	tok.value = val;
	return tok;
}

void Token::print() {
	std::string ty;

	switch (this->type) {
		case Type::T_EOF: ty = "T_EOF"; break;
		case Type::T_EOL: ty = "T_EOL"; break;
		case Type::IGNORE: ty = "IGNORE"; break;
		case Type::STRING: ty = "STRING"; break;
		case Type::NUMBER: ty = "NUMBER"; break;
		case Type::KEYWORD: ty = "KEYWORD"; break;
		case Type::OPERATOR: ty = "OPERATOR"; break;
		case Type::NAME: ty = "NAME"; break;
		case Type::BOOL: ty = "BOOL"; break;
		default: ty = "UNKNOWN"; break;
	}

	std::cout << ty << ": " << this->value << "\n";
}
