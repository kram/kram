// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <unordered_map>
#include <vector>
#include <string>
#include "../vm/instruction.h"
#include "../lexer/token.h"

enum class ON {
	DEFAULT,
	MATH_CONTINUATION,
	PUSH_CLASS,
};

class Parser {
	std::vector<lexer::Token*> tokens;

	int index;
	int lenght;
	bool has_advanced;

	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> comparisions;
	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> startOperators;
	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> leftOnlyInfix;
	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> rightOnlyInfix;

	std::vector<Instruction*> read_file();
	std::vector<Instruction*> read_until_eol();
	
	std::vector<Instruction*> read_until(std::vector<lexer::Token>);
	std::vector<Instruction*> read_until(std::vector<lexer::Token>, std::vector<Instruction*>);

	Instruction* lookahead(Instruction*, ON);
	lexer::Token* get_token();
	lexer::Token* get_and_expect_token(lexer::Token);
	void advance();
	void reverse();
	Instruction* symbol_next(ON);
	Instruction* symbol(lexer::Token*, ON);

	Instruction* oper(lexer::Token*);
	Instruction* oper_list_creation();
	Instruction* oper_list_extraction(Instruction*);

	Instruction* keyword(lexer::Token*);
	Instruction* keyword_class();
	Instruction* keyword_fn();
	Instruction* keyword_if();
	Instruction* keyword_new();
	Instruction* keyword_while();

	Instruction* name(lexer::Token*, ON);
	
	Instruction* boolean(lexer::Token*, ON);
	Instruction* boolean_init(std::string);
	Instruction* boolean_init();
	
	Instruction* number(lexer::Token*, ON);
	Instruction* number_init(std::string);
	Instruction* number_init();

	Instruction* string(lexer::Token*, ON);
	Instruction* string_init(std::string);
	Instruction* string_init();

	Instruction* ignore();
	
	Instruction* assign(Instruction*);
	Instruction* assign_with_type(Instruction*);

	Instruction* math(Instruction*);
	Instruction* push_class(Instruction*);
	Instruction* call(Instruction*, ON);

	public: 
		Parser(std::vector<lexer::Token*>);
		std::vector<Instruction*> run();
		static int infix_priority(lexer::Type);
};