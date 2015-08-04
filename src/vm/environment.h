// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

#include <unordered_map>

#include "value.h"
#include "map.h"

class Environment {
	Kram_Map::map * names;

	public:

		Environment();

		Environment* parent;
		Environment* root;
		bool is_root;

		void set(const char *, Value*, bool root = false);
		//void set_root(const char *, Value*);

		Value* get(const char *, bool root = false);
		//Value* get_root(const char *);

		bool has(const char *);

		Environment* push();
		Environment* pop();
};