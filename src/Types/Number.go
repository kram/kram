package types

import (
	"log"
	"math"
	"strconv"
)

type Number struct {
	Number bool
	Value  float64
}

func (self *Number) Init(str string) {
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		log.Panicf("Can not initialize Number as %s", str)
	}

	self.Value = value
}

func (self Number) Type() string {
	return "Number"
}

func (self *Number) ToString() string {
	return strconv.FormatFloat(self.Value, 'f', 6, 64)
}

func (self *Number) Math(method string, right Type) Type {

	r, ok := right.(*Number)

	if !ok {
		log.Panicf("You can not apply %s to a %s() with a %s()", method, self.Type(), right.Type())
	}

	val := float64(0)

	if method == "+" || method == "-" || method == "*" || method == "/" || method == "%" || method == "**" {
		switch method {
		case "+":
			val = self.Value + r.Value
		case "-":
			val = self.Value - r.Value
		case "*":
			val = self.Value * r.Value
		case "/":
			val = self.Value / r.Value
		case "%":
			val = math.Mod(self.Value, r.Value)
		case "**":
			val = math.Pow(self.Value, r.Value)
		}

		num := Number{}
		num.Value = val
		return &num
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	num := Number{}
	num.Value = val

	return &num
}

func (self *Number) Compare(method string, right Type) Type {

	r, ok := right.(*Number)

	if !ok {
		log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())
	}

	b := false

	switch method {
	case ">":
		b = self.Value > r.Value
	case "<":
		b = self.Value < r.Value
	case ">=":
		b = self.Value >= r.Value
	case "<=":
		b = self.Value <= r.Value
	case "==":
		b = self.Value == r.Value
	case "!=":
		b = self.Value != r.Value
	default:
		log.Panicf("%s() is not implementing %s", self.Type(), method)
	}

	bl := Bool{}
	bl.Value = b

	return &bl
}
