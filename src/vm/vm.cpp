#import "vm.h"

#import <iostream>

#import "library.h"
//#import "libraries/user/class.h"
#import "libraries/user/function.h"
#import "libraries/IO/io.h"

Value VM::assign(Instruction ins) {
	this->names[ins.name] = this->run(ins.right[0]);

	return Value::NUL();
}

Value VM::literal(Instruction ins) {
	// The value is already pre-calcualted by the parser
	return ins.value;
}

Value VM::name(Instruction ins) {

	// TODO
	if (this->names.find(ins.name) == this->names.end()) {
		return Value::STRING(ins.name);	
	}

	std::cout << "name/get\n";

	Value res = this->names[ins.name];

	std::cout << "name/got\n";

	std::cout << res.print() << "\n";

	return res;
}

Value VM::math(Instruction ins) {

	int res = 0;
	int l = this->run(ins.left[0]).getNumber();
	int r = this->run(ins.right[0]).getNumber();

	switch (ins.type) {
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

	return Value::NUMBER(res);
}

Value VM::if_case(Instruction ins) {
	return Value::NUL();
}

Value VM::ignore(Instruction ins) {
	return Value::NUL();
}

Value VM::push_class(Instruction ins) {
	// Get the name of the class to push
	Value name = this->name(ins.left[0]);

	// Add a pointer to the class to the back (aka top) of the stack
	this->lib_stack.push_back(&this->names[name.getString()]);

	// Run the right part
	return this->run(ins.right);
}

Value VM::function(Instruction ins) {
	Function fn;
	fn.REFERENCE("Function");
	fn.init();
	fn.set_content(ins.right);
	// fn.set_vm(this);

	std::cout << "Defined function\n";

	return fn;
}

Value VM::call(Instruction ins) {

	ins.print();

	// Get the method name or function declaration
	Value fun = this->name(ins.left[0]);

	if (fun.type != Type::REFERENCE) {
		return this->call_library(ins);
	}

	std::cout << "Calling function\n" << fun.print();


	if (ins.right.size() == 1) {
		fun.execMethod("exec", this->run(ins.right[0]));
	} else {
		fun.execMethod("exec", Value::NUL());
	}

	// TODO: Return values
	return Value::NUL();
}

Value VM::call_library(Instruction ins) {

	std::cout << "Calling library\n";

	// Get the method name
	Value name = this->name(ins.left[0]);

	// Get the library from the top of the stack
	Value* lib = this->lib_stack[this->lib_stack.size() - 1];

	// Get the first parameter
	// TODO: Allow for more parameters (and none)
	Value params = this->run(ins.right[0]);

	// Call the method
	lib->execMethod(name.getString(), params);

	// TODO: Return values
	return Value::NUL();
}

Value VM::run(Instruction ins) {
	switch (ins.instruction) {
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

	return Value::NUL();
}

Value VM::run(std::vector<Instruction> ins) {

	std::cout << "block()\n";

	for (Instruction i : ins) {
		this->run(i);
	}

	return Value::NUL();
}

void VM::boot(std::vector<Instruction> ins) {
	IO io;
	io.init();
	io.REFERENCE("IO");
	this->names["IO"] = io;

	std::cout << "Booted\n";

	this->run(ins);
}