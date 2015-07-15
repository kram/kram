#ifndef LEXER_TOKEN_H
#define LEXER_TOKEN_H

#include <string>
#include <unordered_map>

namespace lexer {
	enum class Type {
		T_EOF,
		T_EOL,

		IGNORE,

		STRING,
		NUMBER,
		NAME,

		BOOL,
		BOOL_TRUE,
		BOOL_FALSE,

		KEYWORD,
		// KEYWORD_VAR,

		OPERATOR,
	
		OPERATOR_COMMA,         // ,
		OPERATOR_COLON,         // :
		OPERATOR_DOUBLE_COLON,  // :: (OPERATOR_PAAMAYIM_NEKUDOTAYIM)
		OPERATOR_SEMICOLON,     // ;
		OPERATOR_EQ,            // =
		OPERATOR_NOT_EQ,        // !=
		OPERATOR_EQEQ,          // ==
		OPERATOR_GT,            // >
		OPERATOR_GTEQ,          // >=
		OPERATOR_LT,            // <
		OPERATOR_LTEQ,          // <=
		OPERATOR_DOUBLE_AND,    // &&
		OPERATOR_DOUBLE_OR,     // ||
		OPERATOR_PLUS_PLUS,     // ++
		OPERATOR_MINUS_MINUS,   // --
		OPERATOR_PLUS,          // +
		OPERATOR_MINUS,         // -
		OPERATOR_MUL,           // *
		OPERATOR_DIV,           // /
		OPERATOR_DOT,           // .
		OPERATOR_2DOT,          // ..
		OPERATOR_3DOT,          // ...
		OPERATOR_PAREN_L,       // (
		OPERATOR_PAREN_R,       // )
		OPERATOR_COLON_EQ,      // :=
	};

	struct Token {
		public:
			Type type;
			Type sub;
			std::string value;

			Token();

			static Token T_EOF();
			static Token T_EOL();
			static Token IGNORE();
			static Token STRING(std::string);
			static Token NUMBER(std::string);
			static Token KEYWORD(std::string);
			static Token OPERATOR(std::string);
			static Token NAME(std::string);
			static Token BOOL(std::string);

			std::string print();
			static std::string print(Type);
			static Type Trans(std::string);
	};
}

#endif