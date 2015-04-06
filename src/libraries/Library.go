package libraries

import (
	"github.com/zegl/Gus/src/types"
)

type Library struct{}
func (lib *Library) Init(str string) {}
func (lib *Library) Instance() (types.Lib, string) { return &Library{}, lib.Type() }
func (lib *Library) Type() string { return "Library" }
func (lib *Library) ToString() string { return lib.Type() }

func (lib *Library) TypeWithLib(l types.Lib) *types.Type {
	class := types.Type{}
	class.InitWithLib(l)

	return &class
}