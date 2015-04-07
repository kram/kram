// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/zegl/Gus/src/types"
	"log"
	"strings"
)

type Map struct {
	Builtin
	items map[string]*types.Type
	order []string // To keep the order
}

func (self Map) Instance() (types.Lib, string) {
	return &Map{}, "Map"
}

func (self Map) Type() string {
	return "Map"
}

// Map can not be initialized with Init
// see InitWithParams
func (self *Map) Init(str string) {}

func (self *Map) InitWithParams(params []*types.Type) {
	self.items = make(map[string]*types.Type)
	self.order = make([]string, 0)

	is_key := true

	for i, key := range params {
		if is_key {
			k := key.ToString()
			self.order = append(self.order, k)
			self.items[k] = params[i+1]
		}

		is_key = !is_key
	}
}

func (self *Map) Set(params []*types.Type) {
	if len(params) != 2 {
		log.Panic("Library_Map::Set() expected exactly 2 parameters")
	}

	key := params[0].ToString()
	value := params[1]

	// New item
	if _, ok := self.items[key]; !ok {
		self.order = append(self.order, key)
	}

	self.items[key] = value
}

func (self *Map) Get(params []*types.Type) *types.Type {
	if len(params) != 1 {
		log.Panic("Library_Map::Get() expected exactly 1 parameter")
	}

	key := params[0].ToString()

	if res, ok := self.items[key]; ok {
		return res
	}

	log.Panicf("Library_Map::Get() no such key %s", key)

	// Will never be reached
	return self.Null()
}

func (self *Map) ToString() string {
	str := "{\n"

	items := make([]string, 0)

	for _, key := range self.order {
		s := "    "
		s += "\"" + key + "\": "
		s += self.items[key].ToString()

		items = append(items, s)
	}

	str += strings.Join(items, ",\n")

	str += "\n}"

	return str
}

func (self *Map) Length() int {
	return len(self.items)
}
