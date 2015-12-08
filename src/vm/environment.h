// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <unordered_map>
#include <array>
#include "value.h"

class Environment {
	std::unordered_map<std::string, Value*>* names;
	// Value* names_arr[1000];
	std::array<Value*, 1000> names_arr;

	public:

		Environment();

		Environment* parent;
		Environment* root;
		bool is_root;

		void set(const std::string&, Value*);
		void set(size_t, Value*);

		void set_root(std::string, Value*);

		void update(const std::string&, Value*);
		void update(size_t, Value*);

		Value* get(const std::string&);
		Value* get(size_t);

		Value* get_root(const std::string&);

		bool has(const std::string&);
		bool exists(const std::string&);

		Environment* push();
		Environment* pop();
};