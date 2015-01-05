package main

import (
	"log"
)

type String struct {
	String bool
	Value string
}

func (s *String) Init(str string) {
	s.Value = str
}

func (s *String) Math(method string, right Type) Type {
	log.Panicf("String() is not implementing %s", method)

	// This code will never be reached

	res := Bool{}
	res.Init("false")

	return &res
}

func (s *String) toString() string {
	return s.Value
}