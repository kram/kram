#include <iostream>
#include <unordered_map>

class Value;

class Library {
	typedef void (*method)(Value);

	protected:
		std::unordered_map<std::string, method> methods;
		method single_method;

		void add_method(std::string name, method m) {
			std::cout << "add_method() " << name << "\n";
			this->methods[name] = m;
		}

		void add_method(method m) {
			std::cout << "single add_method()\n";
			this->single_method = m;
		}

	public:
		Value call(Value val) {
			std::cout << "call single\n";
			this->single_method(val);
			return Value::NUL();
		}

		Value call(std::string name, Value val) {

			std::cout << "Lib::call() " << name << "\n";

			if (this->methods.find(name) == this->methods.end()) {
				std::cout << "UNKNOWN METHOD: " << name << "\n";
				return Value::NUL();
			}

			std::cout << "Pre\n";

			method m = this->methods[name];

			std::cout << "Post\n";

			m(val);

			return Value::NUL();
		}

		void init(void);

		// Public to be useable from within functions
		VM* vm;
		void set_vm(VM* vm) {
			this->vm = vm;
		}

		std::string print(void) {
			std::string res;

			for(auto kv : this->methods) {
				res += kv.first + "() ";
			}

			return res;
		}
};