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

	// lib_stack contains pointers to values that have been pushed to the stack
	// See: push_class
	std::vector<Value*> lib_stack;

	// environment contains the user scope of names, another layer is added to the environment when a function is executed,
	// and is popped afterwards.
	// See: assign, define_class, env_pop, and env_push
	Environment* environment;

	// Assigns a Value* to the environment, or to the most recently pushed class on lib_stack
	// Set does the same, but can only update old values
	Value* assign(Instruction*, vm::ON);
	Value* set(Instruction*, vm::ON);

	// Literal values are values such as "Hello", 1337, and false.
	// This metod converts them to a Value* that can be used in the VM
	Value* literal(Instruction*);

	// Does a lookup of names that either is defined in the environment or in a library / class
	// Returns either the Value* directly or a Type::NAME (see value.h) that later will be used when executing a method
	Value* name(Instruction*, vm::ON);

	// math() and math_number() does all mathematical operations
	Value* math(Instruction*);
	Value* math_number(Instruction*, Value*, Value*);

	// if_case() takes an expression, if the expression is true the left-side of the instruction is executed, if not the right one is
	Value* if_case(Instruction*);

	// Does nothing. (Returns NUL)
	Value* ignore(Instruction*);

	// Pushes left to the stack and continues to execute right
	Value* push_class(Instruction*);

	// Creates a new function that later can be either directly executed, assigned to the environment, or assigned to a class
	Value* function(Instruction*);

	// Creates a new instance of an already existing Class / Library
	Value* create_instance(Instruction*);

	// Execute right for as long as left executes to true
	Value* loop_while(Instruction*);

	Value* kw_return(Instruction*);

	// Defines a new class on the root level of environment
	// Will execute assignment instructions (other instructions are illegal and will stop the program)
	Value* define_class(Instruction*);

	// Shortcut for creating lists with an optional list of values
	// ["a", "b", "c", "d"]
	// evaluates to what in the longer version is
	// new List().Push("a", "b", "c", "d")
	Value* list_create(Instruction*);

	// Shortcut for extracing a value form a list
	// li[123] instead of li.At(123)
	Value* list_extract(Instruction*);

	Value* list_range(int start, int end, bool inclusive);

	// call(), call_library(), and call_builtin() will each call and execute methods of different types
	Value* call(Instruction*, vm::ON);
	Value* call_library(Instruction*);
	Value* call_builtin(Instruction*, Value*);

	// Steps in and out of the environment
	void env_pop();
	void env_push();

	bool function_should_return();

	Value* value_null;
	Value* value_number [255];
	Value* value_bool [2];

	void init_default_values();
	Value* get_value_null();
	Value* get_value_number(double);
	Value* get_value_bool(bool);

	public:
		// Initialize the VM
		void boot(std::vector<Instruction*>);

		// Run a single or multiple instructions
		Value* run(Instruction*);
		Value* run(Instruction*, vm::ON);
		Value* run(std::vector<Instruction*>);
		std::vector<Value*> run_vector(std::vector<Instruction*>);

		// Set and get values from the environment
		// See environment
		void name_create(const std::string&, Value*);
		void name_update(const std::string&, Value*);
		void name_create_root(std::string, Value*);
		Value* name_get(const std::string&);
		Value* name_get_root(const std::string&);

		// Stack for functions, used together with the return keyword
		void in_function_push();
		void in_function_pop();

	private:
		std::vector<bool> function_return_stack;
};

#endif