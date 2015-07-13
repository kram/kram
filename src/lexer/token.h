#ifndef LEXER_TOKEN_H
#define LEXER_TOKEN_H

#include <string>

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
		KEYWORD_VAR, // var

		OPERATOR,
		OPERATOR_EQ, // =
		OPERATOR_COLON, // ,
		OPERATOR_SEMICOLON, // ;
		OPERATOR_DOT, // .
		OPERATOR_PAREN_L, // (
		OPERATOR_PAREN_R, // )
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
			void print();
	};
}

#endif