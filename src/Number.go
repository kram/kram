package main

import (
	"strconv"
	"log"
)

type Number struct {
	Number bool
	value float64
}

func (n *Number) Init(str string) bool {	
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		log.Panicf("Can not initialize Number as %s", str)
		return false
	}

	n.value = value

	return true
}

func (n *Number) toString() string {
	return strconv.FormatFloat(n.value, 'f', 6, 64)
}

func (n *Number) Add(num float64) float64 {
	n.value += num

	return n.value
}

func (n *Number) Sub(num float64) float64 {
	n.value -= num

	return n.value
}

func (n *Number) Divide(num float64) float64 {
	n.value /= num

	return n.value
}

func (n *Number) Multiply(num float64) float64 {
	n.value *= num

	return n.value
}

func (n *Number) LessThan(num float64) bool {
	if n.value < num {
		return true
	}

	return false
}

func (n *Number) BiggerThan(num float64) bool {
	if n.value > num {
		return true
	}

	return false
}

func (n *Number) EqualTo(num float64) bool {
	if n.value == num {
		return true
	}

	return false
}