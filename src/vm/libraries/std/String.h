// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>

#include "../../../third_party/utf8/utf8.h"
#include "../../value.h"

std::string utf8chr(int cp)
{
    char c[5]={ 0x00,0x00,0x00,0x00,0x00 };
    if     (cp<=0x7F) { c[0] = cp;  }
    else if(cp<=0x7FF) { c[0] = (cp>>6)+192; c[1] = (cp&63)+128; }
    else if(0xd800<=cp && cp<=0xdfff) {} //invalid block of utf8
    else if(cp<=0xFFFF) { c[0] = (cp>>12)+224; c[1]= ((cp>>6)&63)+128; c[2]=(cp&63)+128; }
    else if(cp<=0x10FFFF) { c[0] = (cp>>18)+240; c[1] = ((cp>>12)&63)+128; c[2] = ((cp>>6)&63)+128; c[3]=(cp&63)+128; }
    return std::string(c);
}

class String: public Value {

	/**
	 * String::bytes()
	 *
	 * Used as "test".bytes()
	 *
	 * Returns the amount of chars used to build the whole string.
	 */
	static Value* bytes(Value* self, std::vector<Value*> val) {
		return new Value(Type::NUMBER, val[0]->getString().length());
	}

	/**
	 * String::length()
	 *
	 * Returnes the amount of "glyphs" that the string consists of
	 * Eg, a, ö, and 🌍 all have a length of 1
	 *
	 */
	static Value* length(Value* self, std::vector<Value*> val) {
		std::string full_string = val[0]->getString();
		char* c_str = (char*) full_string.c_str();

		char* str_iterator = c_str;
		char* str_end = c_str + strlen(c_str) + 1;

		int length = utf8::distance(str_iterator, str_end);

		// Do not include NUL
		--length;

		return new Value(Type::NUMBER, length);
	}

	/**
	 * String::at()
	 *
	 *
	 */
	static Value* at(Value* self, std::vector<Value*> val) {
		std::string full_string = val[0]->getString();
		char* c_str = (char*) full_string.c_str();

		char* str_iterator = c_str;
		char* str_end = c_str + strlen(c_str) + 1;

		int length = utf8::distance(str_iterator, str_end);

		// Do not include NUL
		--length;

		int get_char_at = val[1]->getNumber();

		if (get_char_at >= length) {
			std::cout << "String::At() expects a number smaller than the lenght of the string\n";
			exit(0);
		}

		for (int i = 0; i < get_char_at; i++) {
			utf8::next(str_iterator, str_end);
		}

		auto res = utf8::next(str_iterator, str_end);
		return new Value(Type::STRING, utf8chr(res));

		std::cout << "String::At() failed\n";
		exit(0);

		return new Value(Type::NUL);
	}

	public:
		void init() {
			this->add_method("Length", this->length);
			this->add_method("Bytes", this->bytes);
			this->add_method("At", this->at);
		}
};