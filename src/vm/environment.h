// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <unordered_map>
#include "value.h"

class Environment {
	std::unordered_map<std::string, Value*> names;
	std::unordered_map<std::string, Value*> all_names;

	public:

		Environment();

		Environment* parent;
		Environment* root;
		bool is_root;

		void set(std::string, Value*);
		void set_root(std::string, Value*);

		Value* get(std::string);
		Value* get_root(std::string);

		bool has(std::string);

		Environment* push();
		Environment* pop();
};