// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <iostream>

void kr_print_help() {
	std::cout << "Kram " << KR_VERSION << "\n"
	          << "Run a script with `kram path/to/file.kr`\n"
	          << "Available commands:\n"
	          << "\t--help: This text\n"
	          << "\t--version: Print version related information";
}

void kr_print_version() {
	std::cout << "Kram " << KR_VERSION << "\n"
	          << "Built at " << __DATE__ << " " << __TIME__ << "\n";
}