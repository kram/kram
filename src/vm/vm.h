#ifndef VM_H
#define VM_H

#import <unordered_map>
#import <vector>

#import "value.h"
#import "Instruction.h"

class VM {

	std::vector<Value*> lib_stack;
	std::unordered_map<std::string, Value*> names;

	Value* assign(Instruction);
	Value* literal(Instruction);
	Value* name(Instruction);
	Value* math(Instruction);
	Value* if_case(Instruction);
	Value* ignore(Instruction);
	Value* push_class(Instruction);
	Value* function(Instruction);
	Value* call(Instruction);
	Value* call_library(Instruction);

	public:
		void boot(std::vector<Instruction>);

		// Adressable from libraries and what not
		Value* run(Instruction);
		Value* run(std::vector<Instruction>);
		void set_name(std::string, Value*);
};

#endif