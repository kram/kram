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
	return this->read_until(std::vector<lexer::Token>{
		lexer::Token::T_EOF()
	});
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
					return res;
				}
			}
		}

		// Test again
		this->advance();
		lexer::Token next = this->get_token();

		for (lexer::Token unt : until) {
			if (unt.type == next.type && unt.sub == next.sub) {
				return res;
			}
		}

		// We may continue
		res.push_back(this->symbol(next));

		first = false;
	}

	return res;
}

Instruction Parser::lookahead(Instruction prev, ON on) {
	this->advance();

	lexer::Token next = this->get_token();

	// PushClass
	// IO.Println("123")
	//   ^
	if (next.type == lexer::Type::OPERATOR && next.sub == lexer::Type::OPERATOR_DOT) {
		return this->push_class(prev);
	}

	next.print();

	// Call
	// IO.Println("123")
	//           ^
	if (next.type == lexer::Type::OPERATOR && next.sub == lexer::Type::OPERATOR_PAREN_L) {
		return this->call(prev, on);
	}

	this->reverse();

	return prev;
}

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
	if (this->has_advanced) {
		this->index++;
	} else {
		this->has_advanced = true;
	}
}
void Parser::reverse() {
	this->index--;
}

Instruction Parser::symbol_next() {
	this->advance();
	return this->symbol(this->get_token());
}

Instruction Parser::symbol(lexer::Token tok) {
	switch (tok.type) {
		case lexer::Type::KEYWORD: return this->keyword(tok); break;
		case lexer::Type::NUMBER: return this->number(tok); break;
		case lexer::Type::NAME: return this->name(tok); break;
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
	Instruction ins(Ins::ASSIGN);

	this->advance();

	// Get name
	ins.name = this->get_and_expect_token(lexer::Token::NAME("")).value;

	// Expect an equal-sign
	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("="));

	ins.right = this->read_until_eol();

	return ins;
}

//Instruction Parser::keyword_if(lexer::Token);

Instruction Parser::name(lexer::Token tok) {
	Instruction ins(Ins::NAME);
	ins.name = tok.value;

	return this->lookahead(ins, ON::DEFAULT);
}

Instruction Parser::number(lexer::Token tok) {
	Instruction ins(Ins::LITERAL);
	ins.value = Value::NUMBER(std::stoi(tok.value));

	return this->lookahead(ins, ON::DEFAULT);
}

//Instruction Parser::oper(lexer::Token);
//Instruction Parser::ignore(lexer::Token);
//Instruction Parser::bl(lexer::Token);
//Instruction Parser::math(Instruction);

Instruction Parser::push_class(Instruction prev) {
	Instruction ins(Ins::PUSH_CLASS);
	ins.left = std::vector<Instruction>{ prev };
	ins.right = std::vector<Instruction>{ this->symbol_next() };

	ins.print();

	return this->lookahead(ins, ON::PUSH_CLASS);
}

Instruction Parser::call(Instruction prev, ON on) {
	Instruction ins(Ins::CALL);
	ins.left = std::vector<Instruction>{ prev };

	// Read until ) or ,
	ins.right = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR(")"),
		lexer::Token::OPERATOR(","),
	});

	return this->lookahead(ins, ON::DEFAULT);
}