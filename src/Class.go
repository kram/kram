package main

import (
	"log"
)

type Class struct {
	String bool
	Value  string
}

func (class *Class) Init(str string) {
	class.Value = str
}

func (class Class) Type() string {
	return "Class"
}

func (class *Class) Math(method string, right Type) Type {

	log.Panicf("Class() is not implementing %s", method)

	// This code will never be reached

	res := Bool{}
	res.Init("false")

	return &res
}

func (class *Class) toString() string {
	return "Class"
}