#include "parser.h"
#include <iostream>

Parser::Parser(std::vector<lexer::Token> tokens) {
	this->tokens = tokens;
	this->lenght = tokens.size();

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
	// IO::Println("123")
	//   ^^
	// IO.Println("123")
	//   ^
	if (next.type == lexer::Type::OPERATOR && (next.sub == lexer::Type::OPERATOR_DOT || next.sub == lexer::Type::OPERATOR_DOUBLE_COLON)) {
		return this->push_class(prev);
	}

	// Call
	// IO::Println("123")
	//           ^
	if (next.type == lexer::Type::OPERATOR && next.sub == lexer::Type::OPERATOR_PAREN_L) {
		return this->call(prev, on);
	}

	// Assignment
	// num := 100
	//     ^^
	if (next.type == lexer::Type::OPERATOR && next.sub == lexer::Type::OPERATOR_COLON_EQ) {
		return this->assign(prev);
	}

	// Create new variable of type
	// num : Number
	//           ^
	if (next.type == lexer::Type::OPERATOR && next.sub == lexer::Type::OPERATOR_COLON) {
		return this->assign_with_type(prev);
	}

	if (next.type == lexer::Type::OPERATOR) {
		if (this->startOperators.find(next.sub) != this->startOperators.end()) {
			return this->math(prev);
		}
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

Instruction Parser::symbol_next() {
	this->advance();
	return this->symbol(this->get_token());
}

Instruction Parser::symbol(lexer::Token tok) {
	switch (tok.type) {
		case lexer::Type::KEYWORD: return this->keyword(tok); break;
		case lexer::Type::NUMBER: return this->number(tok); break;
		case lexer::Type::NAME: return this->name(tok); break;

		case lexer::Type::T_EOL:
		case lexer::Type::T_EOF:
			return this->ignore();
			break;

		// Just a bit of silencing
		default: break;
	}

	std::cout << "Unknown symbol: " << tok.print() << "\n";
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

Instruction Parser::keyword(lexer::Token tok) {
	switch (tok.sub) {
		// case lexer::Type::KEYWORD_VAR: return this->keyword_var(tok); break;
		case lexer::Type::KEYWORD_CLASS: return this->keyword_class(); break;
		default: break;
	}

	std::cout << "Unknown keyword: " << tok.print() << "\n";
	exit(0);
}

Instruction Parser::keyword_class() {
	Instruction ins(Ins::DEFINE_CLASS);

	this->advance();

	lexer::Token n = this->get_and_expect_token(lexer::Token::NAME(""));
	ins.name = n.value;

	this->advance();
	this->get_and_expect_token(lexer::Token::OPERATOR("{"));

	this->advance();
	ins.right = this->read_until(std::vector<lexer::Token>{
		lexer::Token::OPERATOR("}")
	});

	return ins;
}

Instruction Parser::assign(Instruction prev) {
	Instruction ins(Ins::ASSIGN);

	if (prev.instruction != Ins::NAME) {
		std::cout << ":= expects previous instruction to be of type NAME";
		exit(0);
	}

	// Get name from previous OP
	ins.name = prev.name;
	ins.right = this->read_until_eol();

	return ins;
}

Instruction Parser::assign_with_type(Instruction prev) {
	Instruction ins(Ins::ASSIGN);

	if (prev.instruction != Ins::NAME) {
		std::cout << ": expects previous instruction to be of type NAME";
		exit(0);
	}

	// Get name from previous OP
	ins.name = prev.name;

	// Read the rest of the row
	std::vector<Instruction> row = this->read_until_eol();

	// Expect exactly one name
	if (row.size() != 1 || row[0].instruction != Ins::NAME) {
		std::cout << ": needs to be followed by exactly one NAME (aka the type)";
		exit(0);
	}

	std::string type = row[0].name;

	if (type == "Number") {
		ins.right.push_back(this->number_init());
	} else {
		std::cout << "Does not know how to init type " << type << " after :";
		exit(0);
	}

	return ins;
}

//Instruction Parser::keyword_if(lexer::Token);

Instruction Parser::name(lexer::Token tok) {
	Instruction ins(Ins::NAME);
	ins.name = tok.value;

	return this->lookahead(ins, ON::DEFAULT);
}

Instruction Parser::number(lexer::Token tok) {
	Instruction ins = this->number_init(tok.value);
	return this->lookahead(ins, ON::DEFAULT);
}

Instruction Parser::number_init(std::string val) {
	Instruction ins(Ins::LITERAL);
	ins.value = Value::NUMBER(std::stoi(val));
	return ins;
}

Instruction Parser::number_init() {
	Instruction ins(Ins::LITERAL);
	ins.value = Value::NUMBER(0);
	return ins;
}

//Instruction Parser::oper(lexer::Token);

Instruction Parser::ignore() {
	Instruction res(Ins::IGNORE);
	return res;
}

//Instruction Parser::bl(lexer::Token);

Instruction Parser::math(Instruction prev) {

	// Get the current token
	lexer::Token current = this->get_token();

	// Create a new instruction
	Instruction math(Ins::MATH);

	// Set the mathematical operator, eg + or -
	math.type = current.sub;

	if (prev.instruction == Ins::LITERAL || prev.instruction == Ins::NAME) {
		math.left = std::vector<Instruction> { prev };
		math.right = std::vector<Instruction> { this->symbol_next() };

		// Verify that the ordering (infix_priority()) is correct
		// Left is either a LITERAL or NAME, and right is a (new) MATH
		if (math.right[0].instruction == Ins::MATH) {

			// The ordering is wrong, and we need to correct this
			// [a, *, [b, +, c]] -> [[a, *, b], +, c]
			// This part is a little bit well, confusing and tight. But hey, it is a side-project after all.
			if (Parser::infix_priority(math.type) > Parser::infix_priority(math.right[0].type)) {

				Instruction right = math.right[0];
				Instruction res(Ins::MATH);
				Instruction left(Ins::MATH);

				left.type = math.type;
				left.left = math.left;
				left.right = right.left;
				
				res.left = std::vector<Instruction> { left };
				
				res.type = right.type;
				res.right = right.right;

				return res;
			}
		}
	}

	return math;
}

Instruction Parser::push_class(Instruction prev) {
	Instruction ins(Ins::PUSH_CLASS);
	ins.left = std::vector<Instruction>{ prev };
	ins.right = std::vector<Instruction>{ this->symbol_next() };

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