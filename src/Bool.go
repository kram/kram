package main

type Bool struct {
	Bool  bool
	Value bool
}

func (self *Bool) Init(str string) {

	if str == "true" {
		self.Value = true
	} else {
		self.Value = false
	}
}

func (b Bool) Type() string {
	return "Bool"
}

func (self *Bool) toString() string {

	if self.Value {
		return "true"
	}

	return "false"
}

func (self *Bool) Math(method string, right Type) Type {

	if r, ok := right.(*Bool); ok {
		if self.Value == r.Value {
			bl := Bool{}
			bl.Init("true")

			return &bl
		}
	}

	bl := Bool{}
	bl.Init("false")

	return &bl
}
