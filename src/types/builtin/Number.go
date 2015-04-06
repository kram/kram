package builtin

import (
	"log"
	"math"
	"strconv"
	"../" // types
)

type Number struct {
	Value  float64
}

func (self Number) Instance() (types.Lib, string) {
	return &Number{}, self.Type()
}

func (self Number) Type() string {
	return "Number"
}

func (self *Number) Init(str string) {
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		log.Panicf("Can not initialize Number as %s", str)
	}

	self.Value = value
}

func (self *Number) ToString() string {
	return strconv.FormatFloat(self.Value, 'f', 6, 64)
}

func (self *Number) Math(method string, right *types.Type) *types.Type {

	r, ok := right.Extension.(*Number)
	_, is_null := right.Extension.(*Null)

	if !ok && !is_null {
		log.Panicf("You can not apply %s to a %s() with a %s()", method, self.Type(), right.Type())
	}

	val := float64(0)

	if method == "+" || method == "-" || method == "*" || method == "/" || method == "%" || method == "**" || method == ".." || method == "..." || method == "++" || method == "--" {
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
		case "++":
			self.Value++
			val = self.Value
		case "--":
			self.Value--
			val = self.Value
		case "..", "...":

			list := List{}
			class := types.Type{}
			class.InitWithLib(&list)

			i := self.Value

			for {
				if (method == ".." && i >= r.Value) || (method == "..." && i > r.Value) {
					break
				}

				// Create number object
				num := Number{}
				num.Value = i

				n := types.Type{}
				n.InitWithLib(&num)

				list.Items = append(list.Items, &n)

				i++
			}

			return &class
		}

		num := Number{}
		num.Value = val

		res := types.Type{}
		res.InitWithLib(&num)

		return &res
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	return &types.Type{}
}

func (self *Number) Compare(method string, right *types.Type) *types.Type {

	r, ok := right.Extension.(*Number)

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
	bl.Set(b)

	res := types.Type{}
	res.InitWithLib(&bl)

	return &res
}