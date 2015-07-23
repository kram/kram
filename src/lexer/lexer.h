// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <unordered_map>
#include <vector>
#include "token.h"

namespace lexer {
	class Lexer {
		std::unordered_map<std::string, bool> keywords;
		std::string row;
		size_t index;
		char current;

		Token* next(void);
		char char_at_pos(size_t);

		Token* comment(void);
		Token* name(void);
		Token* number(void);
		Token* string(void);
		Token* oper(void);

		public:
			Lexer();
			std::vector<Token*> parse_file(std::string);
			static void print(std::vector<Token*>);
	};
}