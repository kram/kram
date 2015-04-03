package main

import (
	"log"
	"strings"
)

type Library_Map struct {
	*Library
	items map[string]Type
}

func (self *Library_Map) Instance() (Lib, string) {
	return &Library_Map{}, "Map"
}

func (self *Library_Map) Init(vm *VM, params []Type) {

	self.items = make(map[string]Type)

	is_key := true

	for i, key := range params {
		if is_key {
			self.items[key.ToString()] = params[i + 1]
		}

		is_key = !is_key
	}
}

func (self *Library_Map) Set(vm *VM, params []Type) {
	if len(params) != 2 {
		log.Panic("Library_Map::Set() expected exactly 2 parameters")
	}

	key := params[0].ToString()
	value := params[1]

	self.items[key] = value
}

func (self *Library_Map) ToString() string {
	str := "{\n"

	items := make([]string, 0)

	for key, value := range self.items {
		s := "    "
		s += "\"" + key + "\": "
		s += value.ToString()
		items = append(items, s)
	}

	str += strings.Join(items, ",\n")

	str += "\n}"

	return str
}

func (self *Library_Map) Length() int {
	return len(self.items)
}