// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#ifndef VM_H
#define VM_H

#include <unordered_map>
#include <vector>

#include "value.h"
#include "instruction.h"
#include "environment.h"

class VM {

	std::vector<Value*> lib_stack;
	Environment* environment;

	Value* assign(Instruction*);
	Value* literal(Instruction*);
	Value* name(Instruction*);

	Value* math(Instruction*);
	Value* math_number(Instruction*, Value*, Value*);

	Value* if_case(Instruction*);
	Value* ignore(Instruction*);
	Value* push_class(Instruction*);
	Value* function(Instruction*);
	Value* create_instance(Instruction*);
	Value* loop_while(Instruction*);

	Value* call(Instruction*);
	Value* call_library(Instruction*);
	Value* call_builtin(Instruction*);

	std::vector<Value*> run_vector(std::vector<Instruction*>);

	void env_pop();
	void env_push();

	public:
		void boot(std::vector<Instruction*>);

		// Adressable from libraries and what not
		Value* run(Instruction*);
		Value* run(std::vector<Instruction*>);

		void set_name(std::string, Value*);
		Value* get_name(std::string);
};

#endif