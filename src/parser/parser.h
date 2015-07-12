#include <unordered_map>
#include <vector>

class Parser {
	std::vector<Token> tokens;

	int index;
	int lenght;
	bool has_advanced;

	std::unordered_map<std::string, bool> comparisions;
	std::unordered_map<std::string, bool> startOperators;
	std::unordered_map<std::string, bool> leftOnlyInfix;
	std::unordered_map<std::string, bool> rightOnlyInfix;
};