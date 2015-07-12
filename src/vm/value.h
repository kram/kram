enum class Type {
	NUL,
	STRING,
	NUMBER,
	REFERENCE
};

class Value {
	Type type;
	union {
		std::string string;
		int number;
	};
};