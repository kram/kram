#include <string>

enum class Type {
	T_EOF,
	T_EOL,
	IGNORE,
	STRING,
	NUMBER,
	KEYWORD,
	OPERATOR,
	NAME,
	BOOL
};

struct Token {
	public:
		Type type;
		std::string value;

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