#import <unordered_map>
#import <vector>

#import "value.h"
#import "Instruction.h"

#import "libraries/library.h"
#import "libraries/IO/io.h"

class VM {

	std::unordered_map<std::string, Library> env;
	std::vector<Library*> lib_stack;

	std::unordered_map<std::string, Value> names;

	Value run(Instruction);

	Value assign(Instruction);
	Value literal(Instruction);
	Value name(Instruction);
	Value math(Instruction);
	Value if_case(Instruction);
	Value ignore(Instruction);
	Value push_class(Instruction);
	Value call(Instruction);
	Value block(std::vector<Instruction>);

	public:
		void boot(std::vector<Instruction>);
};