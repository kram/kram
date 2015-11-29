// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#define KR_VERSION "0.0.1"

#include <iostream>

#include "lexer/lexer.h"
#include "parser/parser.h"
#include "vm/instruction.h"
#include "vm/vm.h"
#include "vm/output.h"
#include "cli/helper.h"

Output* kram_output_stream;

int main(int argc, char** argv) {

	bool debug = false;
	std::string file = "";

	// Parse arguments
	// This probable needs to become way more advanced in the future, but it is enough for now...
	if (argc > 1) {
		for (int i = 1; i < argc; i++) {
			std::string arg(argv[i]);

			if (arg.substr(0, 2) == "--") {

				if (arg == "--debug") {
					debug = true;
				} else if (arg == "--help") {
					kr_print_help();
					exit(0);
				} else if (arg == "--version") {
					kr_print_version();
					exit(0);
				}

			} else {
				file = arg;
			}
		}
	} else {
		kr_print_help();
		exit(0);
	}

	// Define output stream
	kram_output_stream = new Output(std::cout);

	lexer::Lexer lexer;
	std::vector<lexer::Token*> tokens = lexer.parse_file(file);

	if (debug) {
		lexer::Lexer::print(tokens);
	}

	Parser parser (tokens);
	std::vector<Instruction*> instructions = parser.run();

	if (debug) {
		std::cout << "Printing ins:\n";

		for (Instruction* ins : instructions) {
			ins->print();
			std::cout << "--------\n";
		}
	}

	VM vm;
	vm.boot(instructions);

	return 0;
}