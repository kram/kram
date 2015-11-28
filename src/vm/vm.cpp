// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include "vm.h"

#include <iostream>

#include "libraries/user/class.h"
#include "libraries/user/function.h"
#include "libraries/IO/io.h"
#include "libraries/std/Number.h"
#include "libraries/std/Map.h"
#include "libraries/std/List.h"
#include "libraries/std/String.h"

Value* VM::assign(Instruction* ins, vm::ON on) {

	if (on == vm::ON::PUSHED_CLASS) {
		Value* top = this->lib_stack.back();

		if (top->type != Type::CLASS) {
			std::cout << "The assign operation (:=) is not allowed on this object (userland classes only)\n";
			exit(0);
		}

		// Convert top to a Class*
		Class* cl = static_cast<Class*>(top);
		cl->set_value(ins->name, this->run(ins->right[0]));
	} else {
		this->name_create(ins->name, this->run(ins->right[0]));
	}

	return new Value(Type::NUL);
}

Value* VM::set(Instruction* ins, vm::ON on) {
	if (on == vm::ON::PUSHED_CLASS) {
		Value* top = this->lib_stack.back();

		if (top->type != Type::CLASS) {
			std::cout << "The set operation (=) is not allowed on this object (userland classes only)\n";
			exit(0);
		}

		// Convert top to a Class*
		Class* cl = static_cast<Class*>(top);
		cl->update_value(ins->name, this->run(ins->right[0]));
	} else {
		this->name_update(ins->name, this->run(ins->right[0]));
	}

	return new Value(Type::NUL);
}

Value* VM::literal(Instruction* ins) {
	// The value is already pre-calcualted by the parser
	return ins->value;
}

Value* VM::name(Instruction* ins, vm::ON on) {

	// Check if the value exists on the pushed class
	if (on == vm::ON::PUSHED_CLASS) {
		Value* top = this->lib_stack.back();

		// Get value from Class directly
		if (top->type == Type::CLASS) {
			Class* cl = static_cast<Class*>(top);
			return cl->get_value(ins->name);
		}

		// Create a Value of type NAME that can be used later
		if (top->type == Type::REFERENCE) {

			if (top->has_method(ins->name)) {
				return new Value(Type::NAME, ins->name);
			}

			std::cout << "REFERENCE: Undefined name: " << ins->name << "\n";
			exit(0);
		}

		// The type is something else, usually a builtin (eg. NUMBER)
		return new Value(Type::NAME, ins->name);
	}

	// Check if the value exists on the stack
	return this->name_get(ins->name);
}

Value* VM::math(Instruction* ins) {
	Value* left = this->run(ins->left[0]);
	Value* right = this->run(ins->right[0]);

	if (left->type != right->type) {
		std::cout << "math() expects equal types, got " << left->print(true) << " and " << right->print(true) << "\n";
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

		case lexer::Type::OPERATOR_2DOT:
		case lexer::Type::OPERATOR_3DOT:
			return this->list_range(l, r, ins->type == lexer::Type::OPERATOR_3DOT);

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

	this->env_push();

	Value* res = this->run(ins->center[0]);

	if (res->type != Type::BOOL) {
		std::cout << "If-case must evaluate to a BOOL\n";
		exit(0);
	}

	// Was true
	if (res->getBool()) {
		auto res = this->run(ins->left);
		this->env_pop();
		return res;
	}

	// Has else-part
	if (ins->right.size() > 0) {
		auto res = this->run(ins->right);
		this->env_pop();
		return res;
	}

	this->env_pop();

	// Return NUL otherwise
	return new Value(Type::NUL);
}

Value* VM::loop_while(Instruction* ins) {

	while (true) {
		Value* res = this->run(ins->left[0]);

		if (res->type != Type::BOOL) {
			std::cout << "While-case must evaluate to a BOOL\n";
			return new Value(Type::NUL);
		}

		// Not true anymore
		if (res->getBool() == false) {
			return new Value(Type::NUL);
		}

		this->env_push();
		this->run(ins->right);
		this->env_pop();
	}
}

Value* VM::kw_return(Instruction* ins) {

	auto res = this->run(ins->right);

	this->function_return_stack.pop_back();
	this->function_return_stack.push_back(true);

	return res;
}

void VM::in_function_push() {
	this->function_return_stack.push_back(false);
}

void VM::in_function_pop() {
	this->function_return_stack.pop_back();
}

bool VM::function_should_return() {
	return this->function_return_stack.back();
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

	if (ins->right.size() != 1) {
		std::cout << "push_class() expects exactly 1 child. This is a bug, please report it! :)\n";
		exit(0);
	}

	// Run the right part
	auto res = this->run(ins->right[0], vm::ON::PUSHED_CLASS);

	this->lib_stack.pop_back();

	return res;
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

Value* VM::define_class(Instruction* ins) {
	Class* cl = new Class();
	cl->set_type(Type::CLASS);
	cl->init();

	cl->vm = this;

	for (Instruction* sub : ins->right) {
		if (sub->instruction != Ins::ASSIGN) {
			std::cout << "Class definitions can only contain assignments\n";
			exit(0);
		}

		cl->set_value(sub->name, this->run(sub->right[0]));
	}

	// Set in the global scope
	this->name_create_root(ins->name, cl);

	return cl;
}

Value* VM::list_create(Instruction* ins) {
	List* list = new List();
	list->set_type(Type::REFERENCE);
	list->init();

	if (ins->right.size() > 0) {
		std::vector<Value*> values;

		for (Instruction* i : ins->right) {
			values.push_back(this->run(i));
		}

		list->push(values);
	}

	return list;
}

Value* VM::list_extract(Instruction* ins) {
	Value* list     = this->run(ins->left[0]);
	Value* position = this->run(ins->right[0]);

	return list->exec_method("At", std::vector<Value *> { position });
}

Value* VM::list_range(int start, int end, bool inclusive) {
	List* list = new List();
	list->set_type(Type::REFERENCE);
	list->init();

	if (!inclusive) {
		end--;
	}

	if (end < start) {
		std::cout << "Range: End needs to be bigger or equal to Start\n";
		exit(0);
	}

	std::vector<Value*> items;

	for (int i = start; i <= end; i++) {
		items.push_back(new Value(Type::NUMBER, i));
	}

	list->push(items);

	return list;
}

Value* VM::create_instance(Instruction* ins) {
	Value* original = this->name_get(ins->name);

	// Kram-defined classes
	if (original->type == Type::CLASS) {
		// Cast as class
		Class* cl = static_cast<Class*>(original);
		return cl->new_instance();
	}

	Value* instance = original->exec_method("new", this->run_vector(ins->right));
	return instance;
}

Value* VM::call(Instruction* ins, vm::ON on) {

	this->env_push();

	// Get the method name or function declaration
	Value* fun = this->name(ins->left[0], on);

	// Result pointer
	Value* res;

	// Function and class methods
	if (fun->type == Type::FUNCTION) {
		// Parse arguments first
		auto arguments = this->run_vector(ins->right);

		// Assign "self" to the parent
		if (on == vm::ON::PUSHED_CLASS) {
			this->name_create("self", this->lib_stack.back());
		}

		// Execute function
		res = fun->exec_method("exec", arguments);

	// Pushed classes (the class has to be fetched from the stack)
	} else if (on == vm::ON::PUSHED_CLASS && fun->type == Type::NAME) {
		Value* top = this->lib_stack.back();

		switch (top->type) {
			case Type::NUMBER:
			case Type::STRING:
			case Type::BOOL:
				res = this->call_builtin(ins, fun);
				break;

			default:
				res = top->exec_method(fun->getString(), this->run_vector(ins->right));
				break;
		}

	// Default action is built in libraries
	} else {
		res = this->call_library(ins);
	}

	this->env_pop();

	return res;
}

Value* VM::call_library(Instruction* ins) {

	// Get the method name
	Value* name = this->name(ins->left[0], vm::ON::DEFAULT);

	// Get the library from the top of the stack
	Value* lib = this->lib_stack.back();

	// Execute the parameters and call the method
	return lib->exec_method(name->getString(), this->run_vector(ins->right));
}

Value* VM::call_builtin(Instruction* ins, Value* name) {

	// Get the value from the top of the stack
	Value* builtin_value = this->lib_stack.back();

	// Get library
	std::string lib_name;// = new std::string;

	switch (builtin_value->type) {
		case Type::NUMBER: lib_name = "Number"; break;
		case Type::STRING: lib_name = "String"; break;
		default:
			std::cout << "call_builtin(): Can not call on " << builtin_value->print() << "\n";
			exit(0);
			break;
	}

	Value* lib = this->name_get(lib_name);

	// Initialize parameter vector with the first value beeing the bultin itself
	std::vector<Value*> params { builtin_value };

	for (auto param : ins->right) {
		params.push_back(this->run(param));
	}

	// Call the method
	return lib->exec_method(name->getString(), params);
}

std::vector<Value*> VM::run_vector(std::vector<Instruction*> instructions) {
	std::vector<Value*> res;

	for (Instruction* i : instructions) {
		res.push_back(this->run(i));
	}

	return res;
}

Value* VM::run(Instruction* ins) {
	return this->run(ins, vm::ON::DEFAULT);
}

Value* VM::run(Instruction* ins, vm::ON on) {
	switch (ins->instruction) {
		case Ins::ASSIGN:          return this->assign(ins, on);      break;
		case Ins::SET:             return this->set(ins, on);         break;
		case Ins::LITERAL:         return this->literal(ins);         break;
		case Ins::NAME:            return this->name(ins, on);        break;
		case Ins::MATH:            return this->math(ins);            break;
		case Ins::IF:              return this->if_case(ins);         break;
		case Ins::IGNORE:          return this->ignore(ins);          break;
		case Ins::PUSH_CLASS:      return this->push_class(ins);      break;
		case Ins::CALL:            return this->call(ins, on);        break;
		case Ins::FUNCTION:        return this->function(ins);        break;
		case Ins::CREATE_INSTANCE: return this->create_instance(ins); break;
		case Ins::WHILE:           return this->loop_while(ins);      break;
		case Ins::DEFINE_CLASS:    return this->define_class(ins);    break;
		case Ins::LIST_CREATE:     return this->list_create(ins);     break;
		case Ins::LIST_EXTRACT:    return this->list_extract(ins);    break;
		case Ins::RETURN:          return this->kw_return(ins);       break;

		case Ins::FUNCTION_PARAMETER:
			std::cout << "Could not evaluate FUNCTION_PARAMETER\n";
			exit(1);
			break;
	}

	return new Value(Type::NUL);
}

Value* VM::run(std::vector<Instruction*> ins) {

	Value* last;

	for (Instruction* i : ins) {

		last = this->run(i, vm::ON::DEFAULT);

		if (i->instruction == Ins::RETURN) {
			return last;
		}

		if (this->function_should_return()) {
			return last;
		}
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

	String* str = new String();
	str->set_type(Type::REFERENCE);
	str->init();
	this->environment->set_root("String", str);

	Map* map = new Map();
	map->set_type(Type::REFERENCE);
	map->init();
	this->environment->set_root("Map", map);

	List* list = new List();
	list->set_type(Type::REFERENCE);
	list->init();
	this->environment->set_root("List", list);

	this->function_return_stack = std::vector<bool>{false};

	this->run(ins);
}