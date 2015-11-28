// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "parser.h"
#include <iostream>

Parser::Parser(std::vector<lexer::Token*> tokens) {

	// Initialize variables
	this->tokens = tokens;
	this->lenght = tokens.size();
	this->index = 0;
	this->has_advanced = false;

	// Hashmap of comparisions
	this->comparisions[lexer::Type::OPERATOR_EQEQ] = true;
	this->comparisions[lexer::Type::OPERATOR_GT] = true;
	this->comparisions[lexer::Type::OPERATOR_GTEQ] = true;
	this->comparisions[lexer::Type::OPERATOR_LT] = true;
	this->comparisions[lexer::Type::OPERATOR_LTEQ] = true;
	this->comparisions[lexer::Type::OPERATOR_DOUBLE_AND] = true;
	this->comparisions[lexer::Type::OPERATOR_DOUBLE_OR] = true;

	// 123++
	this->leftOnlyInfix[lexer::Type::OPERATOR_PLUS_PLUS] = true;
	this->leftOnlyInfix[lexer::Type::OPERATOR_MINUS_MINUS] = true;

	// -123
	this->rightOnlyInfix[lexer::Type::OPERATOR_MINUS] = true;

	// List of all operators starting a new sub-expression
	this->startOperators[lexer::Type::OPERATOR_EQEQ] = true;
	this->startOperators[lexer::Type::OPERATOR_GT] = true;
	this->startOperators[lexer::Type::OPERATOR_GTEQ] = true;
	this->startOperators[lexer::Type::OPERATOR_LT] = true;
	this->startOperators[lexer::Type::OPERATOR_LTEQ] = true;
	this->startOperators[lexer::Type::OPERATOR_DOUBLE_AND] = true;
	this->startOperators[lexer::Type::OPERATOR_DOUBLE_OR] = true;
	this->startOperators[lexer::Type::OPERATOR_PLUS_PLUS] = true;
	this->startOperators[lexer::Type::OPERATOR_MINUS_MINUS] = true;
	this->startOperators[lexer::Type::OPERATOR_MINUS] = true;
	this->startOperators[lexer::Type::OPERATOR_PLUS] = true;
	this->startOperators[lexer::Type::OPERATOR_MUL] = true;
	this->startOperators[lexer::Type::OPERATOR_DIV] = true;
	this->startOperators[lexer::Type::OPERATOR_PAREN_L] = true;
	this->startOperators[lexer::Type::OPERATOR_EQ] = true;
	this->startOperators[lexer::Type::OPERATOR_2DOT] = true;
	this->startOperators[lexer::Type::OPERATOR_3DOT] = true;
}

std::vector<Instruction*> Parser::run() {
	return this->read_file();
}

std::vector<Instruction*> Parser::read_file() {
	return this->read_until(std::vector<lexer::Token>{
		lexer::Token::T_EOF()
	});
}

std::vector<Instruction*> Parser::read_until_eol() {
	return this->read_until(std::vector<lexer::Token>{
		lexer::Token::T_EOL()
	});
}

std::vector<Instruction*> Parser::read_until(std::vector<lexer::Token> until) {
	std::vector<Instruction*> res;
	return this->read_until(until, res, ON::DEFAULT);
}

std::vector<Instruction*> Parser::read_until(std::vector<lexer::Token> until, std::vector<Instruction*> res) {
	return this->read_until(until, res, ON::DEFAULT);
}

std::vector<Instruction*> Parser::read_until(std::vector<lexer::Token> until, std::vector<Instruction*> res, ON on) {
	// Everything can end at a EOF
	until.push_back(lexer::Token::T_EOF());

	bool first = true;

	while (true) {

		// First, test the previous token
		lexer::Token* prev = this->get_token();

		if (!first) {
			for (lexer::Token unt : until) {
				if ((unt.type == lexer::Type::T_EOL || unt.type == lexer::Type::T_EOF || (unt.type == lexer::Type::OPERATOR && unt.sub == lexer::Type::OPERATOR_SEMICOLON)) && prev->type == unt.type) {
					return res;
				}
			}
		}

		// Test again
		this->advance();
		lexer::Token* next = this->get_token();

		for (lexer::Token unt : until) {
			if (unt.type == next->type && unt.sub == next->sub) {
				return res;
			}
		}

		Instruction* sym = this->symbol(next, on);

		if (sym->instruction != Ins::IGNORE) {
			res.push_back(sym);	
		}

		first = false;
	}

	return res;
}

Instruction* Parser::lookahead(Instruction* prev, ON on) {
	this->advance();

	lexer::Token* next = this->get_token();

	// PushClass
	// IO::Println("123")
	//   ^^
	// IO.Println("123")
	//   ^
	if (next->type == lexer::Type::OPERATOR && (next->sub == lexer::Type::OPERATOR_DOT || next->sub == lexer::Type::OPERATOR_DOUBLE_COLON)) {
		return this->push_class(prev);
	}

	// Call
	// IO::Println("123")
	//            ^
	if (next->type == lexer::Type::OPERATOR && next->sub == lexer::Type::OPERATOR_PAREN_L) {
		return this->call(prev, on);
	}

	// Assignment
	// num := 100
	//     ^^
	if (next->type == lexer::Type::OPERATOR && next->sub == lexer::Type::OPERATOR_COLON_EQ) {
		return this->assign(prev);
	}

	// Set
	// num = 100
	//     ^
	//
	// Parameter with default value
	// fn (a = 100)
	//
	if (next->type == lexer::Type::OPERATOR && next->sub == lexer::Type::OPERATOR_EQ) {
		if (on == ON::FUNCTION_PARAMETER_LIST) {
			return this->function_parameter_with_default_value(prev);
		} else {
			return this->set(prev);
		}
	}

	// Create new variable of type
	// num : Number
	//     ^
	if (next->type == lexer::Type::OPERATOR && next->sub == lexer::Type::OPERATOR_COLON) {
		return this->assign_with_type(prev);
	}

	if (on != ON::MATH_CONTINUATION && next->type == lexer::Type::OPERATOR) {
		if (this->startOperators.find(next->sub) != this->startOperators.end()) {
			return this->math(prev);
		}
	}

	// List extraction
	// li[123]
	//   ^
	if (next->type == lexer::Type::OPERATOR && next->sub == lexer::Type::OPERATOR_SQUARE_PAREN_LEFT) {
		return this->oper_list_extraction(prev);
	}

	this->reverse();

	return prev;
}

lexer::Token* Parser::get_token() {
	return this->tokens[this->index];
}

lexer::Token* Parser::get_and_expect_token(lexer::Token expect) {
	lexer::Token* tok = this->get_token();

	if (expect.type != tok->type) {
		std::cout << "Expected: " << expect.print() << "\n";
		std::cout << "Got: " << tok->print() << "\n";
		exit(0);
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

Instruction* Parser::symbol_next(ON on) {
	this->advance();
	return this->symbol(this->get_token(), on);
}

Instruction* Parser::symbol(lexer::Token* tok, ON on) {
	switch (tok->type) {

		// Keywords have yet another switch, go there.
		case lexer::Type::KEYWORD: return this->keyword(tok); break;

		// Operators also have another switch, go there.
		case lexer::Type::OPERATOR: return this->oper(tok); break;

		// Literals
		case lexer::Type::NUMBER: return this->number(tok, on); break;
		case lexer::Type::STRING: return this->string(tok, on); break;
		case lexer::Type::NAME: return this->name(tok, on); break;
		case lexer::Type::BOOL: return this->boolean(tok, on); break;

		case lexer::Type::T_EOL:
		case lexer::Type::T_EOF:
			return this->ignore();
			break;

		// Just a bit of silencing
		default: break;
	}

	std::cout << "Unknown symbol: " << tok->print() << "\n";
	exit(0);
}

int Parser::infix_priority(lexer::Type in) {
	switch (in) {
		// && ||
		case lexer::Type::OPERATOR_DOUBLE_AND:
		case lexer::Type::OPERATOR_DOUBLE_OR:
			return 30;
			break;

		// Comparisions
		case lexer::Type::OPERATOR_EQEQ:
		case lexer::Type::OPERATOR_NOT_EQ:
		case lexer::Type::OPERATOR_LT:
		case lexer::Type::OPERATOR_LTEQ:
		case lexer::Type::OPERATOR_GT:
		case lexer::Type::OPERATOR_GTEQ:
			return 40;
			break;

		case lexer::Type::OPERATOR_PLUS:
		case lexer::Type::OPERATOR_MINUS:
			return 50;
			break;

		case lexer::Type::OPERATOR_MUL:
		case lexer::Type::OPERATOR_DIV:
			return 60;
			break;

		case lexer::Type::OPERATOR_2DOT:
		case lexer::Type::OPERATOR_3DOT:
			return 70;
			break;

		case lexer::Type::OPERATOR_DOT:
		case lexer::Type::OPERATOR_PAREN_L:
		case lexer::Type::OPERATOR_EQ:
		case lexer::Type::OPERATOR_PLUS_PLUS:
		case lexer::Type::OPERATOR_MINUS_MINUS:
			return 80;
			break;

		default:
			return 0;
			break;
	}

	return 0;
}

Instruction* Parser::oper(lexer::Token* tok) {
	switch (tok->sub) {
		case lexer::Type::OPERATOR_SQUARE_PAREN_LEFT: return this->oper_list_creation(); break;
		default: break;
	}

	std::cout << "Unknown operator: " << tok->print() << "\n";
	exit(0);
}

Instruction* Parser::oper_list_creation() {
	Instruction* ins = new Instruction(Ins::LIST_CREATE);

	// Get all peices (seperated by commas)
	do {
		ins->right = this->read_until(std::vector<lexer::Token>{
			lexer::Token::OPERATOR("]"),
			lexer::Token::OPERATOR(","),
		}, ins->right);
	} while (this->get_token()->sub == lexer::Type::OPERATOR_COMMA);

	return ins;
}

Instruction* Parser::oper_list_extraction(Instruction* prev) {
	Instruction* ins = new Instruction(Ins::LIST_EXTRACT);

	/*if (prev->size() != 1) {
		std::cout << "oper_list_extraction: There can only be one list to extract from\n";
		exit(0);
	}*/

	ins->left = std::vector<Instruction*>{ prev };
	ins->right = this->read_until(std::vector<lexer::Token>{ lexer::Token::OPERATOR("]"), });

	if (ins->right.size() != 1) {
		std::cout << "oper_list_extraction: Only one value can be extracted\n";
		exit(0);
	}

	return ins;
}

Instruction* Parser::keyword(lexer::Token* tok) {
	switch (tok->sub) {
		case lexer::Type::KEYWORD_CLASS: return this->keyword_class(); break;
		case lexer::Type::KEYWORD_FN: return this->keyword_fn(); break;
		case lexer::Type::KEYWORD_IF: return this->keyword_if(); break;
		case lexer::Type::KEYWORD_NEW: return this->keyword_new(); break;
		case lexer::Type::KEYWORD_WHILE: return this->keyword_while(); break;
		case lexer::Type::KEYWORD_RETURN: return this->keyword_return(); break;
		default: break;
	}

	std::cout << "Unknown keyword: " << tok->print() << "\n";
	exit(0);
}

Instruction* Parser::keyword_class() {
	Instruction* ins = new Instruction(Ins::DEFINE_CLASS);

	this->advance();

	lexer::Token* n = this->get_and_expect_token(lexer::Token::NAME(""));
	ins->name = n->value;

	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("{"));

	this->advance();
	ins->right = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("}")
	});

	return ins;
}

Instruction* Parser::keyword_if() {
	Instruction* ins = new Instruction(Ins::IF);

	ins->center = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("{")
	});

	this->advance();
	ins->left = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("}")
	});

	// Test if there is any "else" part
	this->advance();
	if (this->get_token()->sub != lexer::Type::KEYWORD_ELSE) {
		this->reverse();
		return ins;
	}

	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("{"));

	ins->right = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("}")
	});

	return ins;
}

Instruction* Parser::keyword_new() {
	Instruction* ins = new Instruction(Ins::CREATE_INSTANCE);

	this->advance();

	// The name of the type/class to push
	lexer::Token* n = this->get_and_expect_token(lexer::Token::NAME(""));
	ins->name = n->value;

	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("("));

	// All parameters
	do {
		ins->right = this->read_until(std::vector<lexer::Token>{
			lexer::Token::OPERATOR(")"),
			lexer::Token::OPERATOR(","),
		}, ins->right);
	} while (this->get_token()->sub == lexer::Type::OPERATOR_COMMA);

	return this->lookahead(ins, ON::DEFAULT);
}

Instruction* Parser::keyword_while() {
	Instruction* ins = new Instruction(Ins::WHILE);

	ins->left = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("{")
	});

	this->advance();
	ins->right = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("}")
	});

	return ins;
}

Instruction* Parser::keyword_return() {
	Instruction* ins = new Instruction(Ins::RETURN);
	ins->right = this->read_until_eol();
	return ins;
}

Instruction* Parser::keyword_fn() {
	Instruction* ins = new Instruction(Ins::FUNCTION);

	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("("));

	// All parameters
	do {
		ins->left = this->read_until(std::vector<lexer::Token>{
			lexer::Token::OPERATOR(")"),
			lexer::Token::OPERATOR(","),
		}, ins->left, ON::FUNCTION_PARAMETER_LIST);
	} while (this->get_token()->sub == lexer::Type::OPERATOR_COMMA);

	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("{"));

	this->advance();
	ins->right = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("}")
	});

	return ins;
}

Instruction* Parser::function_parameter_with_default_value(Instruction* ins) {

	if (ins->instruction != Ins::FUNCTION_PARAMETER) {
		std::cout << "Parameter default value expected Ins::FUNCTION_PARAMETER before =, got:\n";
		ins->print();
		exit(1);
	}

	ins->right = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR(")"),
		lexer::Token::OPERATOR(","),
	});

	// Reverse one step because Parser::keyword_fn() is consuming the same tokens
	this->reverse();

	return ins;
}

Instruction* Parser::set(Instruction* prev) {
	Instruction* ins = new Instruction(Ins::SET);

	if (prev->instruction != Ins::NAME) {
		std::cout << "= expects previous instruction to be of type NAME";
		exit(0);
	}

	// Get name from previous OP
	ins->name = prev->name;
	ins->right = this->read_until_eol();

	return ins;
}

Instruction* Parser::assign(Instruction* prev) {
	Instruction* ins = new Instruction(Ins::ASSIGN);

	if (prev->instruction != Ins::NAME) {
		std::cout << ":= expects previous instruction to be of type NAME";
		exit(0);
	}

	// Get name from previous OP
	ins->name = prev->name;
	ins->right = this->read_until_eol();

	return ins;
}

Instruction* Parser::assign_with_type(Instruction* prev) {
	Instruction* ins = new Instruction(Ins::ASSIGN);

	if (prev->instruction != Ins::NAME) {
		std::cout << ": expects previous instruction to be of type NAME";
		exit(0);
	}

	// Get name from previous OP
	ins->name = prev->name;

	// Read the rest of the row
	std::vector<Instruction*> row = this->read_until_eol();

	// Expect exactly one name
	if (row.size() != 1 || row[0]->instruction != Ins::NAME) {
		std::cout << ": needs to be followed by exactly one NAME (aka the type)";
		exit(0);
	}

	std::string type = row[0]->name;

	if (type == "Number") {
		ins->right.push_back(this->number_init());
	} else if (type == "String") {
		ins->right.push_back(this->string_init());
	} else if (type == "Bool") {
		ins->right.push_back(this->boolean_init());
	} else {
		std::cout << "Does not know how to init type " << type << " after :";
		exit(0);
	}

	return ins;
}

Instruction* Parser::name(lexer::Token* tok, ON on) {

	Instruction* ins;

	if (on == ON::FUNCTION_PARAMETER_LIST) {
		ins = new Instruction(Ins::FUNCTION_PARAMETER);
	} else {
		ins = new Instruction(Ins::NAME);
	}
	
	ins->name = tok->value;

	return this->lookahead(ins, on);
}

Instruction* Parser::boolean(lexer::Token* tok, ON on) {
	Instruction* ins = this->boolean_init(tok->value);
	return this->lookahead(ins, on);
}

Instruction* Parser::boolean_init(std::string val) {
	Instruction* ins = new Instruction(Ins::LITERAL);

	if (val == "true") {
		ins->value = new Value(Type::BOOL, 1);
	} else {
		ins->value = new Value(Type::BOOL, 0);
	}

	return ins;
}

Instruction* Parser::boolean_init() {
	Instruction* ins = new Instruction(Ins::LITERAL);
	ins->value = new Value(Type::BOOL, 0);
	return ins;
}

Instruction* Parser::number(lexer::Token* tok, ON on) {
	Instruction* ins = this->number_init(tok->value);
	return this->lookahead(ins, on);
}

Instruction* Parser::number_init(std::string val) {
	Instruction* ins = new Instruction(Ins::LITERAL);
	ins->value = new Value(Type::NUMBER, std::stod(val));
	return ins;
}

Instruction* Parser::number_init() {
	Instruction* ins = new Instruction(Ins::LITERAL);
	ins->value = new Value(Type::NUMBER, 0);
	return ins;
}

Instruction* Parser::string(lexer::Token* tok, ON on) {
	Instruction* ins = this->string_init(tok->value);
	return this->lookahead(ins, on);
}

Instruction* Parser::string_init(std::string val) {
	Instruction* ins = new Instruction(Ins::LITERAL);
	ins->value = new Value(Type::STRING, val);
	return ins;
}

Instruction* Parser::string_init() {
	Instruction* ins = new Instruction(Ins::LITERAL);
	ins->value = new Value(Type::STRING, "");
	return ins;
}

//Instruction* Parser::oper(lexer::Token);

Instruction* Parser::ignore() {
	return new Instruction(Ins::IGNORE);
}

//Instruction* Parser::bl(lexer::Token);

Instruction* Parser::math(Instruction* prev) {

	// Get the current token
	lexer::Token* current = this->get_token();

	// Create a new instruction
	Instruction* math = new Instruction(Ins::MATH);

	// Set the mathematical operator, eg + or -
	math->type = current->sub;

	math->left = std::vector<Instruction*> { prev };
	math->right = std::vector<Instruction*> { this->symbol_next(ON::MATH_CONTINUATION) };

	// Verify that the ordering (infix_priority()) is correct
	if (prev->instruction == Ins::MATH) {

		// The ordering is wrong, and we need to correct this
		// [a, *, [b, +, c]] -> [[a, *, b], +, c]
		// This part is a little bit well, confusing and tight. But hey, it is a side-project after all.
		if (Parser::infix_priority(math->type) > Parser::infix_priority(prev->type)) {

			Instruction* res = new Instruction(Ins::MATH);
			Instruction* right = new Instruction(Ins::MATH);

			right->left = prev->right;
			right->right = math->right;

			res->left = prev->left;
			res->right = std::vector<Instruction*> { right };

			res->type = prev->type;
			right->type = math->type;

			return this->lookahead(res, ON::DEFAULT);
		}
	}

	return this->lookahead(math, ON::DEFAULT);
}

Instruction* Parser::push_class(Instruction* prev) {
	Instruction* ins = new Instruction(Ins::PUSH_CLASS);
	ins->left = std::vector<Instruction*>{ prev };
	ins->right = std::vector<Instruction*>{ this->symbol_next(ON::DEFAULT) };

	return this->lookahead(ins, ON::PUSH_CLASS);
}

Instruction* Parser::call(Instruction* prev, ON on) {
	Instruction* ins = new Instruction(Ins::CALL);
	ins->left = std::vector<Instruction*>{ prev };

	// All parameters
	do {
		ins->right = this->read_until(std::vector<lexer::Token>{
			lexer::Token::OPERATOR(")"),
			lexer::Token::OPERATOR(","),
		}, ins->right);
	} while (this->get_token()->sub == lexer::Type::OPERATOR_COMMA);

	return this->lookahead(ins, ON::DEFAULT);
}