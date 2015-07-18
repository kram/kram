#include "vm.h"

#include <iostream>

//#include "libraries/user/class.h"
#include "libraries/user/function.h"
#include "libraries/IO/io.h"
#include "libraries/std/Number.h"

void VM::set_name(std::string name, Value* val) {
	this->names[name] = val;
}

Value* VM::get_name(std::string name) {
	if (this->names.find(name) == this->names.end()) {
		std::cout << "Unknown name: " << name << "\n";
		exit(0);
	}

	return this->names[name];
}

Value* VM::assign(Instruction* ins) {
	this->names[ins->name] = this->run(ins->right[0]);

	return new Value(Type::NUL);
}

Value* VM::literal(Instruction* ins) {
	// The value is already pre-calcualted by the parser
	Value* val = new Value(ins->value.type, ins->value.getNumber());

	return val;
}

Value* VM::name(Instruction* ins) {

	// TODO
	if (this->names.find(ins->name) == this->names.end()) {
		return new Value(Type::STRING, ins->name);
	}

	return this->get_name(ins->name);
}

Value* VM::math(Instruction* ins) {

	int res = 0;
	int l = this->run(ins->left[0])->getNumber();
	int r = this->run(ins->right[0])->getNumber();

	switch (ins->type) {
		case lexer::Type::OPERATOR_PLUS:
			res = l + r;
			break;
		case lexer::Type::OPERATOR_MINUS:
			res = l - r;
			break;
		case lexer::Type::OPERATOR_DIV:
			res = l / r;
			break;
		case lexer::Type::OPERATOR_MUL:
			res = l * r;
			break;

		// Ssssh!
		default: break;
	}

	return new Value(Type::NUMBER, res);
}

Value* VM::if_case(Instruction* ins) {
	return new Value(Type::NUL);
}

Value* VM::ignore(Instruction* ins) {
	return new Value(Type::NUL);
}

Value* VM::push_class(Instruction* ins) {
	// Run what we should push
	// Can be a name (name.method() or name::method()),
	// a literal (100.Sqrt())
	// or something else
	Value* push = this->run(ins->left[0]);

	// Add a pointer to the class to the back (aka top) of the stack
	this->lib_stack.push_back(push);

	// Run the right part
	return this->run(ins->right);
}

Value* VM::function(Instruction* ins) {
	Function* fn = new Function();
	fn->type = Type::REFERENCE;
	fn->init();
	
	fn->set_parameters(ins->left);
	fn->set_content(ins->right);

	fn->vm = this;

	return fn;
}

Value* VM::call(Instruction* ins) {

	// Get the method name or function declaration
	Value* fun = this->name(ins->left[0]);

	if (fun->type != Type::REFERENCE) {
		return this->call_library(ins);
	}

	if (ins->right.size() == 1) {
		return fun->execMethod("exec", std::vector<Value*>{ this->run(ins->right[0]) });
	}

	return fun->execMethod("exec", std::vector<Value*>{ new Value() });
}

Value* VM::call_library(Instruction* ins) {

	// Get the method name
	Value* name = this->name(ins->left[0]);

	// Get the library from the top of the stack
	Value* lib = this->lib_stack[this->lib_stack.size() - 1];

	if (lib->type == Type::NUMBER) {
		return this->call_builtin(ins);
	}

	// Get the first parameter
	// TODO: Allow for more parameters (and none)
	Value* params = this->run(ins->right[0]);

	// Call the method
	return lib->execMethod(name->getString(), std::vector<Value*>{ params });
}

Value* VM::call_builtin(Instruction* ins) {

	// Get the method name
	Value* name = this->name(ins->left[0]);

	// Get the value from the top of the stack
	Value* builtin_value = this->lib_stack[this->lib_stack.size() - 1];

	// Get library
	Value* lib;

	switch (builtin_value->type) {
		case Type::NUMBER:
			lib = this->get_name("Number");
			break;
		default:
			std::cout << "call_builtin(): Can not call on " << builtin_value->print() << "\n";
			exit(0);
			break;
	}

	// TODO: Parameters

	// Call the method
	return lib->execMethod(name->getString(), std::vector<Value*>{ builtin_value });
}

Value* VM::run(Instruction* ins) {
	switch (ins->instruction) {
		case Ins::ASSIGN:     return this->assign(ins);     break;
		case Ins::LITERAL:    return this->literal(ins);    break;
		case Ins::NAME:       return this->name(ins);       break;
		case Ins::MATH:       return this->math(ins);       break;
		case Ins::IF:         return this->if_case(ins);    break;
		case Ins::IGNORE:     return this->ignore(ins);     break;
		case Ins::PUSH_CLASS: return this->push_class(ins); break;
		case Ins::CALL:       return this->call(ins);       break;
		case Ins::FUNCTION:   return this->function(ins);   break;
		default: std::cout << "Unknown instruction";        break;
	}

	return new Value(Type::NUL);
}

Value* VM::run(std::vector<Instruction*> ins) {

	Value* last;

	for (Instruction* i : ins) {
		last = this->run(i);
	}

	return last;
}

void VM::boot(std::vector<Instruction*> ins) {
	IO* io = new IO();
	io->type = Type::REFERENCE;
	io->init();
	this->names["IO"] = io;

	Number* number = new Number();
	number->type = Type::REFERENCE;
	number->init();
	this->names["Number"] = number;

	this->run(ins);
}