// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

#include <iostream>

//#include "libraries/user/class.h"
#include "libraries/user/function.h"
#include "libraries/IO/io.h"
#include "libraries/std/Number.h"
#include "libraries/std/Map.h"

void VM::set_name(std::string name, Value* val) {
	this->environment->set(name, val);
}

Value* VM::get_name(std::string name) {
	return this->environment->get(name);
}

void VM::env_pop() {
	this->environment = this->environment->pop();
}

void VM::env_push() {
	this->environment = this->environment->push();
}

Value* VM::assign(Instruction* ins) {
	this->set_name(ins->name, this->run(ins->right[0]));

	return new Value(Type::NUL);
}

Value* VM::literal(Instruction* ins) {
	// The value is already pre-calcualted by the parser
	return ins->value;
}

Value* VM::name(Instruction* ins) {

	// TODO: This method needs to be reworked to work properly...

	if (this->environment->has(ins->name)) {
		return this->get_name(ins->name);
	}

	return new Value(Type::STRING, ins->name);
}

Value* VM::math(Instruction* ins) {
	Value* left = this->run(ins->left[0]);
	Value* right = this->run(ins->right[0]);

	if (left->type != right->type) {
		std::cout << "math() Can not do math on " << left->print() << " and " << right->print() << "\n";
		exit(0);
	}
	
	switch (left->type) {
		case Type::NUMBER:
			return this->math_number(ins, left, right);
			break;

		// Silence the compiler
		default: break;
	}

	std::cout << "math() Does not know how to handle " << left->print() << "\n";
	exit(0);

	return new Value(Type::NUL);
}

Value* VM::math_number(Instruction* ins, Value* left, Value* right) {

	int res_int = 0;
	bool res_bool = false;
	bool is_bool = false;
	int l = left->getNumber();
	int r = right->getNumber();

	switch (ins->type) {
		case lexer::Type::OPERATOR_PLUS:
			res_int = l + r;
			break;

		case lexer::Type::OPERATOR_MINUS:
			res_int = l - r;
			break;

		case lexer::Type::OPERATOR_DIV:
			res_int = l / r;
			break;

		case lexer::Type::OPERATOR_MUL:
			res_int = l * r;
			break;

		case lexer::Type::OPERATOR_LT:
			res_bool = l < r;
			is_bool = true;
			break;

		case lexer::Type::OPERATOR_GTEQ:
			res_bool = l >= r;
			is_bool = true;
			break;

		case lexer::Type::OPERATOR_GT:
			res_bool = l > r;
			is_bool = true;
			break;

		case lexer::Type::OPERATOR_LTEQ:
			res_bool = l <= r;
			is_bool = true;
			break;

		case lexer::Type::OPERATOR_EQEQ:
			res_bool = l == r;
			is_bool = true;
			break;

		case lexer::Type::OPERATOR_NOT_EQ:
			res_bool = l != r;
			is_bool = true;
			break;

		default:
			std::cout << "Unknown math_number() operator\n";
			exit(0);
			break;
	}

	if (is_bool) {
		res_int = (res_bool ? 1 : 0);
		return new Value(Type::BOOL, res_int);
	}

	return new Value(Type::NUMBER, res_int);
}

Value* VM::if_case(Instruction* ins) {

	Value* res = this->run(ins->center[0]);

	if (res->type != Type::BOOL) {
		std::cout << "If-case must evaluate to a BOOL\n";
		exit(0);
	}

	// Was true
	if (res->getBool()) {
		return this->run(ins->left);
	}

	// Has else-part
	if (ins->right.size() > 0) {
		return this->run(ins->right);
	}

	// Return NUL otherwise
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
	fn->set_type(Type::FUNCTION);
	fn->init();
	
	fn->set_parameters(ins->left);
	fn->set_content(ins->right);

	fn->vm = this;

	return fn;
}

Value* VM::create_instance(Instruction* ins) {
	Value* original = this->get_name(ins->name);
	Value* instance = original->execMethod("new", this->run_vector(ins->right));
	return instance;
}

Value* VM::call(Instruction* ins) {

	this->env_push();

	// Get the method name or function declaration
	Value* fun = this->name(ins->left[0]);

	Value* res;

	if (fun->type != Type::REFERENCE && fun->type != Type::FUNCTION) {
		res = this->call_library(ins);
	} else {
		res = fun->execMethod("exec", this->run_vector(ins->right));
	}

	this->env_pop();

	return res;
}

Value* VM::call_library(Instruction* ins) {

	// Get the method name
	Value* name = this->name(ins->left[0]);

	// Get the library from the top of the stack
	Value* lib = this->lib_stack[this->lib_stack.size() - 1];

	if (lib->type == Type::NUMBER) {
		return this->call_builtin(ins);
	}

	// Execute the parameters and call the method
	return lib->execMethod(name->getString(), this->run_vector(ins->right));
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

std::vector<Value*> VM::run_vector(std::vector<Instruction*> instructions) {
	std::vector<Value*> res;

	for (Instruction* i : instructions) {
		res.push_back(this->run(i));
	}

	return res;
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
		case Ins::CREATE_INSTANCE: return this->create_instance(ins); break;
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

	// Create environment
	this->environment = new Environment();
	this->environment->is_root = true;

	IO* io = new IO();
	io->set_type(Type::REFERENCE);
	io->init();
	this->environment->set_root("IO", io);

	Number* number = new Number();
	number->set_type(Type::REFERENCE);
	number->init();
	this->environment->set_root("Number", number);

	Map* map = new Map();
	map->set_type(Type::REFERENCE);
	map->init();
	this->environment->set_root("Map", map);

	this->run(ins);
}