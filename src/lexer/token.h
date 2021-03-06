// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

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
		KEYWORD_CLASS,
		KEYWORD_FN,
		KEYWORD_IF,
		KEYWORD_ELSE,
		KEYWORD_NEW,
		KEYWORD_WHILE,
		KEYWORD_RETURN,

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
		OPERATOR_PAREN_L,       // (
		OPERATOR_PAREN_R,       // )
		OPERATOR_CURLYPAREN_L,  // {
		OPERATOR_CURLYPAREN_R,  // }
		OPERATOR_COLON_EQ,      // :=
		OPERATOR_NOT,           // !

		// Ranges
		OPERATOR_2DOT,          // ..
		OPERATOR_3DOT,          // ...

		// List operations
		OPERATOR_SQUARE_PAREN_LEFT,  // [
		OPERATOR_SQUARE_PAREN_RIGHT, // ]
	};

	// Used to hash lexer::Type for useage in std::unordered_map
	struct EnumClassHash {
		template <typename T>
		std::size_t operator()(T t) const
		{
			return static_cast<std::size_t>(t);
		}
	};

	struct Token {
		public:
			Type type;
			Type sub;
			std::string value;

			Token(Type);
			Token(Type, std::string);

			// TODO: Get rid of the need of these
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