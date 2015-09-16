// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

// This file is based on the works by Chris Wu
// http://chriswu.me/blog/writing-hello-world-in-fcgi-with-c-plus-plus/

#include <stdlib.h>
#include <iostream>
#include <string>

#include "fcgio.h"
#include "fcgi_kram.h"
#include "../vm/output.h"

// Maximum bytes
const unsigned long STDIN_MAX = 1000000;

Output* kram_output_stream;

/**
 * Note this is not thread safe due to the static allocation of the
 * content_buffer.
 */
std::string get_request_content(const FCGX_Request & request) {
	char * content_length_str = FCGX_GetParam("CONTENT_LENGTH", request.envp);
	unsigned long content_length = STDIN_MAX;

	if (content_length_str) {
		content_length = strtol(content_length_str, &content_length_str, 10);
		if (*content_length_str) {
			std::cerr << "Can't Parse 'CONTENT_LENGTH='"
				 << FCGX_GetParam("CONTENT_LENGTH", request.envp)
				 << "'. Consuming stdin up to " << STDIN_MAX << std::endl;
		}

		if (content_length > STDIN_MAX) {
			content_length = STDIN_MAX;
		}
	} else {
		// Do not read from stdin if CONTENT_LENGTH is missing
		content_length = 0;
	}

	char * content_buffer = new char[content_length];
	std::cin.read(content_buffer, content_length);
	content_length = std::cin.gcount();

	// Chew up any remaining stdin - this shouldn't be necessary
	// but is because mod_fastcgi doesn't handle it correctly.

	// ignore() doesn't set the eof bit in some versions of glibc++
	// so use gcount() instead of eof()...
	do std::cin.ignore(1024); while (std::cin.gcount() == 1024);

	std::string content(content_buffer, content_length);
	delete [] content_buffer;
	return content;
}

int main(void) {

	FCGX_Request request;

	FCGX_Init();
	FCGX_InitRequest(&request, 0, 0);

	while (FCGX_Accept_r(&request) == 0) {
		fcgi_streambuf cin_fcgi_streambuf(request.in);
		fcgi_streambuf cout_fcgi_streambuf(request.out);
		fcgi_streambuf cerr_fcgi_streambuf(request.err);

		//std::cin.rdbuf(&cin_fcgi_streambuf);
		//std::cout.rdbuf(&cout_fcgi_streambuf);
		//std::cerr.rdbuf(&cerr_fcgi_streambuf);

		const char * uri = FCGX_GetParam("REQUEST_URI", request.envp);
		const char * script_filename = FCGX_GetParam("SCRIPT_FILENAME", request.envp);

		// TODO: Pass this to kram
		std::string content = get_request_content(request);

		// Simple debug
		std::cout << "Request: " << uri << "\n";
		
		// Redirect output to fcgi
		std::ostream fout(&cout_fcgi_streambuf);
		kram_output_stream = new Output(fout);

		// Print basic headers
		fout << "Content-type: text/html\r\n"
			 << "\r\n";

		run_file(script_filename);

		fout << "\n";

		// Note: the fcgi_streambuf destructor will auto flush
	}

	return 0;
}
