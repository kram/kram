package libraries

import (
	"io/ioutil"
	"log"
	"github.com/zegl/Gus/src/types"
	"github.com/zegl/Gus/src/types/builtin"
)

type Library_File struct {
	*Library
}

func (self *Library_File) Instance() (types.Lib, string) {
	return &Library_File{}, "File"
}

// File.Read()
// @param path String
// @return String
func (self Library_File) Read(params []*types.Type) *types.Type {

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

	str := &builtin.String{}
	str.Init(string(dat))

	return self.TypeWithLib(str)
}

// File.Write()
// @param path String
// @param content String
// @return Bool
func (self Library_File) Write(params []*types.Type) *types.Type {

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

	bl := &builtin.Bool{}
	bl.Init("true")

	return self.TypeWithLib(bl)
}