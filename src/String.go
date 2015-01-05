package main

import (
	"log"
)

type String struct {
	String bool
	Value  string
}

func (s *String) Init(str string) {
	s.Value = str
}

func (s String) Type() string {
	return "String"
}

func (s *String) Math(method string, right Type) Type {

	r, ok := right.(*String)

	if !ok {
		log.Panicf("You can not %s a String with %s", method, right)
	}

	// String concatenation
	if method == "+" {
		str := String{}
		str.Init(s.Value + r.Value)

		return &str
	}

	log.Panicf("String() is not implementing %s", method)

	// This code will never be reached

	res := Bool{}
	res.Init("false")

	return &res
}

func (s *String) toString() string {
	return s.Value
}
