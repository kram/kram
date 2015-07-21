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