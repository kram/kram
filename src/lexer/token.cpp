// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>
#include "token.h"

using namespace lexer;

Token::Token() {
	type = Type::IGNORE;
	sub = Type::IGNORE;
}

Token::Token(Type t) {
	type = t;
	sub = Type::IGNORE;
}

Token::Token(Type t, std::string val) {
	type = t;
	sub = Token::Trans(val);
	value = val;
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
	tok.sub = Token::Trans(val);

	return tok;
}

Token Token::OPERATOR(std::string val) {
	Token tok;
	tok.type = Type::OPERATOR;
	tok.value = val;
	tok.sub = Token::Trans(val);

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

namespace lexer {
	std::unordered_map<std::string, Type> opTrans;
	bool built_op_trans;
}

std::string Token::print() {
	// make sure that opTrans is built
	Token::Trans("");

	std::string res = "";

	res += "T: " + Token::print(this->type) + ", ";
	res += "S: " + Token::print(this->sub) + ", ";
	res += "V: " + this->value;

	return res;
}

std::string Token::print(Type in) {
	for (const auto& kv : opTrans) {
		if (kv.second == in) {
			return kv.first;
		}
	}

	return "UNKNOWN";
}

Type Token::Trans(std::string from) {
	if (!built_op_trans) {

		opTrans["OPERATOR"] = Type::OPERATOR;
		opTrans[","] = Type::OPERATOR_COMMA;
		opTrans[":="] = Type::OPERATOR_COLON_EQ;
		opTrans[":"] = Type::OPERATOR_COLON;
		opTrans["::"] = Type::OPERATOR_DOUBLE_COLON;
		opTrans[";"] = Type::OPERATOR_SEMICOLON;
		opTrans["="] = Type::OPERATOR_EQ;
		opTrans["=="] = Type::OPERATOR_EQEQ;
		opTrans[">"] = Type::OPERATOR_GT;
		opTrans[">="] = Type::OPERATOR_GTEQ;
		opTrans["<"] = Type::OPERATOR_LT;
		opTrans["<="] = Type::OPERATOR_LTEQ;
		opTrans["&&"] = Type::OPERATOR_DOUBLE_AND;
		opTrans["||"] = Type::OPERATOR_DOUBLE_OR;
		opTrans["++"] = Type::OPERATOR_PLUS_PLUS;
		opTrans["--"] = Type::OPERATOR_MINUS_MINUS;
		opTrans["+"] = Type::OPERATOR_PLUS;
		opTrans["-"] = Type::OPERATOR_MINUS;
		opTrans["*"] = Type::OPERATOR_MUL;
		opTrans["/"] = Type::OPERATOR_DIV;
		opTrans["."] = Type::OPERATOR_DOT;
		opTrans[".."] = Type::OPERATOR_2DOT;
		opTrans["..."] = Type::OPERATOR_3DOT;
		opTrans["("] = Type::OPERATOR_PAREN_L;
		opTrans[")"] = Type::OPERATOR_PAREN_R;
		opTrans["{"] = Type::OPERATOR_CURLYPAREN_L;
		opTrans["}"] = Type::OPERATOR_CURLYPAREN_R;

		opTrans["T_EOF"] = Type::T_EOF;
		opTrans["T_EOL"] = Type::T_EOL;
		opTrans["IGNORE"] = Type::IGNORE;
		opTrans["STRING"] = Type::STRING;
		opTrans["NUMBER"] = Type::NUMBER;
		opTrans["NAME"] = Type::NAME;

		opTrans["BOOL"] = Type::BOOL;
		opTrans["true"] = Type::BOOL_TRUE;
		opTrans["false"] = Type::BOOL_FALSE;

		opTrans["KEYWORD"] = Type::KEYWORD;
		opTrans["class"] = Type::KEYWORD_CLASS;
		opTrans["fn"] = Type::KEYWORD_FN;
		opTrans["if"] = Type::KEYWORD_IF;
		opTrans["else"] = Type::KEYWORD_ELSE;
		opTrans["new"] = Type::KEYWORD_NEW;

		built_op_trans = true;
	}

	if (opTrans.find(from) == opTrans.end()) {
		return Type::IGNORE;	
	}

	return opTrans[from];
}