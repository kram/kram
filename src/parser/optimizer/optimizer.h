// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "../../vm/instruction.h"
#include <vector>

class Optimizer {

	size_t depth;
	std::vector<size_t> next_num;
	std::vector<std::unordered_map<std::string, size_t>> names_map;

	void variable_alloc_setup();
	std::vector<Instruction*> variable_alloc_level(std::vector<Instruction*>);
	stack_and_pos variable_alloc_resolve_name(std::string);

	void push();
	void pop();

	Optimizer();

	public:
		static std::vector<Instruction*> variable_alloc(std::vector<Instruction*>);
};