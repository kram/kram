#include <iostream>
#include "src/third_party/utf8/utf8.h"

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


int main(int argc, char** argv) {
	std::string str = "abc🌍abc";

	char* c_str = (char*)str.c_str();    // utf-8 string
	char* str_i = c_str;                 // string iterator
	char* end = c_str+strlen(c_str)+1;

	int length = utf8::distance(str_i, end);

	std::cout << "Lenght: " << length << "\n";

	for (int i = 0; i < length; i++) {
		auto val = utf8::next(str_i, end);
		std::cout << val << "\n";

		std::cout << utf8chr(val) << "\n";
	}
}

