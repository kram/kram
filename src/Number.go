package main

import (
	"strconv"
	"log"
)

type Number struct {
	Number bool
	Value float64
}

func (n *Number) Init(str string) {	
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		log.Panicf("Can not initialize Number as %s", str)
	}

	n.Value = value
}

func (n *Number) toString() string {
	return strconv.FormatFloat(n.Value, 'f', 6, 64)
}

func (n *Number) Math(method string, right Type) Type {
	
	r, ok := right.(*Number)

	if !ok {
		log.Panicf("You can not %s a Number with %s", method, right)
	}

	val := float64(0)

	switch method {
	case "+":
		val = n.Value + r.Value
	case "-":
		val = n.Value - r.Value
	case "*":
		val = n.Value * r.Value
	case "/":
		val = n.Value / r.Value
	default:
		log.Panicf("Number has no such method, %s", method)
	}

	res := Number{}
	res.Value = val

	return &res
}

func (n *Number) LessThan(num float64) bool {
	if n.Value < num {
		return true
	}

	return false
}

func (n *Number) BiggerThan(num float64) bool {
	if n.Value > num {
		return true
	}

	return false
}

func (n *Number) EqualTo(num float64) bool {
	if n.Value == num {
		return true
	}

	return false
}