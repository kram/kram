#include <iostream>
#include <string>
#include <sstream>
#include <unordered_map>
#include <vector>

class VM;

enum Type : int {
	NUL,
	STRING,
	NUMBER,
	BOOL,
	REFERENCE,
	FUNCTION,
	CLASS,
	NAME,
};

class Value {

public: 
	typedef Value* (*Method)(Value*, std::vector<Value*>);
	typedef std::unordered_map<std::string, Method> Methods;

	protected:
		union {
			double number;
			std::string* strval;
			Methods* methods;
			Method single_method;
		} data;

	public:

		Type type;
};

int main()
{
	std::cout << "double: " << sizeof(double) << "\n";
	std::cout << "std::string*: " << sizeof(std::string*) << "\n";
	std::cout << "std::string: " << sizeof(std::string) << "\n";
	std::cout << "Methods*: " << sizeof(Value::Methods*) << "\n";
	std::cout << "Methods: " << sizeof(Value::Methods) << "\n";
	std::cout << "Method: " << sizeof(Value::Method) << "\n";
	std::cout << "int: " << sizeof(int) << "\n";
	std::cout << "Value: " << sizeof(Value) << "\n";

	return 0;
}