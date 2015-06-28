// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"github.com/kram/kram/src/types"
	"io/ioutil"
	"log"
)

type Library_File struct {
	Library
}

func (self Library_File) Instance() (types.Lib, string) { return &Library_File{}, self.Type() }
func (self Library_File) Type() string                  { return "File" }
func (self Library_File) M_Type() *types.Class          { return self.String(self.Type()) }

// File.Read()
// @param path String
// @return String
func (self Library_File) M_Read(params []*types.Class) *types.Class {

	if len(params) != 1 {
		log.Panic("File.Read() expects exactly 1 parameter")
	}

	par := params[0]

	if par.Type() != "String" {
		log.Panic("File.Read() expects parameter 1 to be of type String")
	}

	dat, err := ioutil.ReadFile(par.ToString())

	if err != nil {
		log.Panicf("File.Read(), the file %s was not found", par.ToString())
	}

	return self.String(string(dat))
}

// File.Write()
// @param path String
// @param content String
// @return Bool
func (self Library_File) M_Write(params []*types.Class) *types.Class {

	if len(params) != 2 {
		log.Panic("File.Write() expects exactly 2 parameters")
	}

	for key, param := range params {
		if param.Type() != "String" {
			log.Panic("File.Write() expects parameter %d to be of type String", key)
		}
	}

	path := params[0].ToString()
	content := params[1].ToString()

	data := []byte(content)
	err := ioutil.WriteFile(path, data, 0644)

	if err != nil {
		log.Panicf("File.Write(), could not write to file, %s", path)
	}

	return self.Bool(true)
}
