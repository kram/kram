#ifndef LIBRARY_H
#define LIBRARY_H

#include <unordered_map>
#include <iostream>
#include "../value.h"

class Library {

	typedef void (*method)(Value);

	protected:
		std::unordered_map<std::string, method> methods;

		void add_method(std::string name, method m) {
			std::cout << "add_method(): " << name << "\n";
			this->methods[name] = m;
		}

	public:
		Value call(std::string name, Value val) {
			if (this->methods.find(name) == this->methods.end()) {
				std::cout << "UNKNOWN METHOD: " << name << "\n";
				return Value::NUL();
			}

			method m = this->methods[name];

			m(val);

			return Value::NUL();
		}

		void init(void);
};

#endif