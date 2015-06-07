// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/zegl/Gus/src/types"
	"log"
	"math"
	"strconv"
)

type Number struct {
	Builtin
	Value float64
}

func (self Number) Instance() (types.Lib, string) { return &Number{}, self.Type() }
func (self Number) Type() string                  { return "Number" }
func (self Number) M_Type() *types.Class          { return self.String(self.Type()) }

func (self *Number) Init(str string) {
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		log.Panicf("Can not initialize Number as %s", str)
	}

	self.Value = value
}

func (self *Number) ToString() string {
	return strconv.FormatFloat(self.Value, 'f', -1, 64)
}

func (self *Number) Math(method string, right *types.Class) *types.Class {

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
			if is_null {
				val = -self.Value
			} else {
				val = self.Value - r.Value
			}
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
			class := types.Class{}
			class.InitWithLib(&list)

			i := self.Value

			for {
				if (method == ".." && i >= r.Value) || (method == "..." && i > r.Value) {
					break
				}

				// Create number object
				num := Number{}
				num.Value = i

				n := types.Class{}
				n.InitWithLib(&num)

				list.Items = append(list.Items, &n)

				i++
			}

			return &class
		}

		num := Number{}
		num.Value = val

		res := types.Class{}
		res.InitWithLib(&num)

		return &res
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	return &types.Class{}
}

func (self *Number) Compare(method string, right *types.Class) *types.Class {

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

	res := types.Class{}
	res.InitWithLib(&bl)

	return &res
}

func (self *Number) M_Sqrt(input []*types.Class) *types.Class {
	return self.getNumber(math.Sqrt(self.Value))
}

func (self Number) getNumber(val float64) *types.Class {
	nb := Number{}
	nb.Value = val

	res := types.Class{}
	res.InitWithLib(&nb)

	return &res
}
