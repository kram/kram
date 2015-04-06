package builtin

import (
	"github.com/zegl/Gus/src/types"
)

type Bool struct {
	value bool
}

func (self Bool) Instance() (types.Lib, string) {
	return &Bool{}, self.Type()
}

func (self Bool) Type() string {
	return "Bool"
}

func (self *Bool) Init(str string) {
	if str == "true" {
		self.value = true
	} else {
		self.value = false
	}
}

func (self *Bool) Set(bl bool) {
	self.value = bl
}

func (self *Bool) ToString() string {
	if self.value {
		return "true"
	}

	return "false"
}

func (self *Bool) IsTrue() bool {
	return self.value
}

/*

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
*/