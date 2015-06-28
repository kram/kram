// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"encoding/json"
	"github.com/kram/kram/src/types"
	"github.com/kram/kram/src/types/builtin"
	"log"
)

type Library_JSON struct {
	Library
}

func (self Library_JSON) Instance() (types.Lib, string) { return &Library_JSON{}, self.Type() }
func (self Library_JSON) Type() string                  { return "JSON" }
func (self Library_JSON) M_Type() *types.Class          { return self.String(self.Type()) }

func (self Library_JSON) Error(err error) {
	if err != nil {
		log.Panicf("JSON encountered an error: %s", err)
	}
}

func (self Library_JSON) M_Decode(params []*types.Class) *types.Class {
	if len(params) != 1 {
		log.Panic("JSON::Decode() expects exactly 1 parameter")
	}

	data := []byte(params[0].ToString())

	var obj *json.RawMessage
	err := json.Unmarshal(data, &obj)

	self.Error(err)

	return self.Decode(obj)
}

func (self Library_JSON) Decode(raw *json.RawMessage) *types.Class {

	if raw == nil {
		return self.Null()
	}

	j, err := raw.MarshalJSON()

	// Test if map
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(j, &objmap)

	if err == nil {
		m := self.InitMap()

		for k, v := range objmap {
			m.Set(self.String(k), self.Decode(v))
		}

		return self.fromLib(m)
	}

	// Test if string
	var objstr string
	err = json.Unmarshal(j, &objstr)

	if err == nil {
		return self.String(objstr)
	}

	// Test if number
	var ojbnumber float64
	err = json.Unmarshal(j, &ojbnumber)

	if err == nil {
		return self.Number(ojbnumber)
	}

	// Test if list
	objlist := make([]*json.RawMessage, 0)
	err = json.Unmarshal(j, &objlist)

	if err == nil {
		l := self.InitList()

		for _, v := range objlist {
			value := self.Decode(v)
			l.Push(value)
		}

		return self.fromLib(l)
	}

	// Test if bool
	var objbool bool
	err = json.Unmarshal(j, &objbool)

	if err == nil {
		return self.Bool(objbool)
	}

	log.Panicf("JSON::Decode() could not parse JSON, %s", err)

	return self.Null()
}

func (self Library_JSON) M_Encode(params []*types.Class) *types.Class {
	if len(params) != 1 {
		log.Panic("JSON::Decode() expects exactly 1 parameter")
	}

	return self.Encode(params[0])
}

func (self Library_JSON) Encode(in *types.Class) *types.Class {
	return self.String(string(self.EncodeRaw(in)))
}

func (self Library_JSON) EncodeRaw(in *types.Class) []byte {
	t := in.Type()

	var inter interface{}

	switch t {
	case "Map":
		resmap := make(map[string]*json.RawMessage)
		objmap := in.Extension.(*builtin.Map).GetMap()

		for k, v := range objmap {
			raw := json.RawMessage(self.EncodeRaw(v))
			resmap[k] = &raw
		}

		inter = resmap

	case "List":
		reslist := make([]*json.RawMessage, 0)

		objlist := in.Extension.(*builtin.List).Items

		for _, v := range objlist {
			raw := json.RawMessage(self.EncodeRaw(v))
			reslist = append(reslist, &raw)
		}

		inter = reslist

	case "Number":
		inter = in.Extension.(*builtin.Number).Value
	case "String":
		inter = in.Extension.(*builtin.String).Value
	case "Bool":
		inter = in.Extension.(*builtin.Bool).Value
	case "Null":
		inter = nil
	default:
		log.Panicf("JSON can not encode %f", t)
	}

	res, err := json.Marshal(inter)

	self.Error(err)

	return res
}
