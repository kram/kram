// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "../lexer/lexer.h"
#include "../parser/parser.h"
#include "../vm/instruction.h"
#include "../vm/vm.h"

void run_file(std::string file)
{
	lexer::Lexer lexer;
	auto tokens = lexer.parse_file(file);

	Parser parser (tokens);
	auto instructions = parser.run();

	VM vm;
	vm.boot(instructions);
}