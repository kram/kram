#include <unordered_map>
#include <vector>
#include "token.h"

namespace lexer {
	class Lexer {
		std::unordered_map<std::string, bool> keywords;
		std::string row;
		int index;
		char current;

		Token next(void);
		char char_at_pos(int);

		Token comment(void);
		Token name(void);
		Token number(void);
		Token string(void);
		Token oper(void);

		public:
			Lexer();
			std::vector<Token> parse_file(std::string);
			static void print(std::vector<Token>);
	};
}