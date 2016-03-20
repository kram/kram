#include <unordered_map>
#include <iostream>

class Value {
	public:
		double number;

		Value(double number); 
		static Value* Get(double number); 
};

std::unordered_map<double, Value*> number_vals;

Value::Value(double number)
{
	std::cout << "Created: " << number << "\n";
	this->number = number;
}

Value* Value::Get(double number)
{
	if (number_vals.find(number) == number_vals.end()) {
		Value* res = new Value(number);
		number_vals[number] = res;
		return res;
	}

	std::cout << "Got: " << number << "\n";

	return number_vals[number];
}

int main() {

	auto v1 = Value::Get(10);
	auto v2 = Value::Get(20);
	auto v3 = Value::Get(15);
	auto v4 = Value::Get(10);
	Value::Get(10);
	Value::Get(10);
	Value::Get(10);
}