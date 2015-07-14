#include <iostream>
#include "../library.h"

class IO: public Library {
	static void println(Value val) {
		std::cout << val.print() << "\n";
	}

	public:
		void init() {
			this->add_method("Println", this->println);
		}
};