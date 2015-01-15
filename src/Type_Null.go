package main

import (
	"log"
)

type Null struct {
	Null bool
}

func (self *Null) Init(str string) {}

func (self Null) Type() string {
	return "Null"
}

func (self *Null) ToString() string {
	return "null"
}

func (self *Null) Math(method string, right Type) Type {

	log.Panicf("You can not apply %s to a %s() with a %s()", method, self.Type(), right.Type())

	// Will never be reached

	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (self *Null) Compare(method string, right Type) Type {

	log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())

	// Will never be reached

	bl := Bool{}
	bl.Init("false")

	return &bl
}
