#include <iostream>
#include "../library.h"

class IO: public Library {
	static void println(Value val) {
		std::cout << "PRINT!\n";
		val.print();
	}

	public:
		void init() {
			this->add_method("Println", this->println);
		}
};