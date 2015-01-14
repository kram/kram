package gus

import (
	"log"
)

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

func (self *Bool) ToString() string {

	if self.Value {
		return "true"
	}

	return "false"
}

func (self *Bool) Math(method string, right Type) Type {

	log.Panicf("You can not apply %s to a %s() with a %s()", method, self.Type(), right.Type())

	// Will never be reached

	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (self *Bool) Compare(method string, right Type) Type {

	r, ok := right.(*Bool)

	if !ok {
		log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())
	}

	if self.Value == r.Value {
		bl := Bool{}
		bl.Init("true")

		return &bl
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// Will never be reached

	bl := Bool{}
	bl.Init("false")

	return &bl
}
