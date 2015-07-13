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

	std::unordered_map<std::string, bool> comparisions;
	std::unordered_map<std::string, bool> startOperators;
	std::unordered_map<std::string, bool> leftOnlyInfix;
	std::unordered_map<std::string, bool> rightOnlyInfix;

	std::vector<Instruction> read_file();
	std::vector<Instruction> read_until_eol();
	std::vector<Instruction> read_until(std::vector<lexer::Token>);
	Instruction lookahead(Instruction, ON);
	lexer::Token get_token();
	lexer::Token get_and_expect_token(lexer::Token);
	void advance();
	void reverse();
	Instruction symbol_next();
	Instruction symbol(lexer::Token);
	//uint infix_priority(std::string);
	Instruction keyword(lexer::Token);
	Instruction keyword_var(lexer::Token);
	//Instruction keyword_if(lexer::Token);
	Instruction name(lexer::Token);
	Instruction number(lexer::Token);
	//Instruction oper(lexer::Token);
	//Instruction ignore(lexer::Token);
	//Instruction bl(lexer::Token);
	//Instruction math(Instruction);
	Instruction push_class(Instruction);
	Instruction call(Instruction, ON);

	public: 
		Parser(std::vector<lexer::Token>);
		std::vector<Instruction> run();
};