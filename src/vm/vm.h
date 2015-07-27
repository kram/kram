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

namespace vm {
	enum class ON {
		DEFAULT,
		PUSHED_CLASS,
	};
}

class VM {

	std::vector<Value*> lib_stack;
	Environment* environment;

	Value* assign(Instruction*, vm::ON);
	Value* literal(Instruction*);
	Value* name(Instruction*, vm::ON);

	Value* math(Instruction*);
	Value* math_number(Instruction*, Value*, Value*);

	Value* if_case(Instruction*);
	Value* ignore(Instruction*);
	Value* push_class(Instruction*);
	Value* function(Instruction*);
	Value* create_instance(Instruction*);
	Value* loop_while(Instruction*);

	Value* define_class(Instruction*);

	Value* call(Instruction*, vm::ON);
	Value* call_library(Instruction*);
	Value* call_builtin(Instruction*, Value*);

	std::vector<Value*> run_vector(std::vector<Instruction*>);

	void env_pop();
	void env_push();

	public:
		void boot(std::vector<Instruction*>);

		// Adressable from libraries and what not
		Value* run(Instruction*);
		Value* run(Instruction*, vm::ON);
		Value* run(std::vector<Instruction*>);

		void set_name(std::string, Value*);
		void set_name_root(std::string, Value*);
		Value* get_name(std::string);
		Value* get_name_root(std::string);
};

#endif