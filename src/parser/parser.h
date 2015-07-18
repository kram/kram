#include <unordered_map>
#include <vector>
#include <string>
#include "../vm/instruction.h"
#include "../lexer/token.h"

enum class ON {
	DEFAULT,
	CLASS_BODY,
	PUSH_CLASS,
	METHOD_PARAMETERS,
	ARGUMENTS
};

class Parser {
	std::vector<lexer::Token> tokens;

	int index;
	int lenght;
	bool has_advanced;

	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> comparisions;
	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> startOperators;
	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> leftOnlyInfix;
	std::unordered_map<lexer::Type, bool, lexer::EnumClassHash> rightOnlyInfix;

	std::vector<Instruction> read_file();
	std::vector<Instruction> read_until_eol();
	
	std::vector<Instruction> read_until(std::vector<lexer::Token>);
	std::vector<Instruction> read_until(std::vector<lexer::Token>, std::vector<Instruction>);

	Instruction lookahead(Instruction, ON);
	lexer::Token get_token();
	lexer::Token get_and_expect_token(lexer::Token);
	void advance();
	void reverse();
	Instruction symbol_next();
	Instruction symbol(lexer::Token);

	Instruction keyword(lexer::Token);
	Instruction keyword_class();
	Instruction keyword_fn();
	//Instruction keyword_if();

	Instruction name(lexer::Token);
	
	Instruction number(lexer::Token);
	Instruction number_init(std::string);
	Instruction number_init();

	//Instruction oper(lexer::Token);
	Instruction ignore();
	//Instruction bl(lexer::Token);
	
	Instruction assign(Instruction);
	Instruction assign_with_type(Instruction);

	Instruction math(Instruction);
	Instruction push_class(Instruction);
	Instruction call(Instruction, ON);

	public: 
		Parser(std::vector<lexer::Token>);
		std::vector<Instruction> run();
		static int infix_priority(lexer::Type);
};