// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#ifndef VM_VALUE_H
#define VM_VALUE_H

#include <string>
#include <sstream>
#include <unordered_map>
#include <vector>

#include "map.h"

class VM;

enum class Type {
	NUL,
	STRING,
	NUMBER,
	BOOL,
	REFERENCE,
	FUNCTION,
	CLASS,
	NAME,
};

class Value {

	typedef Value* (*Method)(Value*, std::vector<Value*>);
	// typedef std::unordered_map<std::string, Method> Methods;

	typedef std::unordered_map<const char*, Method, Kram_Map::Hasher, Kram_Map::Equals<const char*> > Methods;

	protected:
		union {
			double number;
			std::string* strval;
			char* charval;
			Methods* methods;
			Method single_method;
		} data;

	public:

		Type type;

		Value();
		Value(Type);
		Value(Type, std::string);
		Value(Type, double);

		void set_type(Type);

		std::string print(bool print_type = false);

		Value* clone();

		std::string getString();
		const char * getCharArray();
		double getNumber();
		bool getBool();

		// Overwritten by references
		void init(void) {}

		// #justlibrarythings
		Value* exec_method(const char *, std::vector<Value*>);
		void add_method(const char *, Method);
		bool has_method(const char *);
};

#endif