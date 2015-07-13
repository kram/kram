#include "parser.h"
#include <iostream>

Parser::Parser(std::vector<lexer::Token> tokens) {
	this->tokens = tokens;
	this->lenght = tokens.size();
}

std::vector<Instruction> Parser::run() {
	return this->read_file();
}

std::vector<Instruction> Parser::read_file() {
	std::vector<Instruction> ins;

	while (true) {

		std::cout << "Loop: " << ins.size() << "\n";

		lexer::Token tok = this->get_token();

		if (tok.type == lexer::Type::T_EOF) {
			break;
		}

		ins.push_back(this->symbol(tok));

		this->advance();
	}

	return ins;
}

std::vector<Instruction> Parser::read_until_eol() {
	return this->read_until(std::vector<lexer::Token>{
		lexer::Token::T_EOL()
	});
}

std::vector<Instruction> Parser::read_until(std::vector<lexer::Token> until) {
	std::vector<Instruction> res;

	// Everything can end at a EOF
	until.push_back(lexer::Token::T_EOF());

	bool first = false;

	while (true) {

		// First, test the previous token
		lexer::Token prev = this->get_token();

		if (!first) {
			for (lexer::Token unt : until) {
				if ((unt.type == lexer::Type::T_EOL || unt.type == lexer::Type::T_EOF || (unt.type == lexer::Type::OPERATOR && unt.sub == lexer::Type::OPERATOR_SEMICOLON)) && prev.type == unt.type) {
					std::cout << "Stopped early\n";
					return res;
				}
			}
		}

		// Test again
		this->advance();
		lexer::Token next = this->get_token();

		std::cout << "read_until() HAS: ";
		next.print();

		for (lexer::Token unt : until) {
			if (unt.type == next.type && unt.sub == next.sub) {
				std::cout << "Stopped reading until\n";
				return res;
			}
		}

		// We may continue
		res.push_back(this->symbol(next));

		first = false;
	}

	return res;
}

//Instruction Parser::lookahead(ON);

lexer::Token Parser::get_token() {
	return this->tokens[this->index];
}

lexer::Token Parser::get_and_expect_token(lexer::Token expect) {
	lexer::Token tok = this->get_token();

	if (expect.type != tok.type) {
		std::cout << "Expected:\n";
		expect.print();
		std::cout << "Got:\n";
		tok.print();
		exit(1);
	}

	return tok;
}

void Parser::advance() {
	this->index++;
}
void Parser::reverse() {
	this->index--;
}

//Instruction Parser::symbol_next();

Instruction Parser::symbol(lexer::Token tok) {
	switch (tok.type) {
		case lexer::Type::KEYWORD: return this->keyword(tok); break;
		case lexer::Type::NUMBER: return this->number(tok); break;
		default: std::cout << "Unknown symbol\n"; tok.print(); break;
	}
}

//uint Parser::infix_priority(std::string);

Instruction Parser::keyword(lexer::Token tok) {
	switch (tok.sub) {
		case lexer::Type::KEYWORD_VAR: return this->keyword_var(tok); break;
		default: std::cout << "Unknown keyword\n"; tok.print(); break;
	}
}

Instruction Parser::keyword_var(lexer::Token tok) {

	std::cout << "keyword_var()\n";

	Instruction ins(Ins::ASSIGN);

	this->advance();

	// Get name
	ins.name = this->get_and_expect_token(lexer::Token::NAME("")).value;

	// Expect an equal-sign
	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("="));

	//this->advance();

	std::cout << "keyword_var() read_until_eol\n";

	ins.right = this->read_until_eol();

	std::cout << "keyword_var() done\n";

	ins.print();

	return ins;
}

//Instruction Parser::keyword_if(lexer::Token);

Instruction Parser::name(lexer::Token tok) {

}

Instruction Parser::number(lexer::Token tok) {
	std::cout << "number()\n";

	Instruction ins(Ins::LITERAL);
	ins.value = Value::NUMBER(std::stoi(tok.value));

	ins.print();

	return ins;
}

//Instruction Parser::oper(lexer::Token);
//Instruction Parser::ignore(lexer::Token);
//Instruction Parser::bl(lexer::Token);
//Instruction Parser::math(Instruction);
//Instruction Parser::push_class(Instruction);
//Instruction Parser::call(Instruction, ON);