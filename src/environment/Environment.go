package environment

import (
	"github.com/zegl/Gus/src/types"
)

type Environment struct {
	Env       map[string]*types.Type
	HasParent bool
	Parent    *Environment
}

func (env *Environment) Init() {
	env.Env = make(map[string]*types.Type)
}

func (env *Environment) Pop() *Environment {

	// No parent
	if !env.HasParent {
		return &Environment{}
	}

	return env.Parent
}

func (env *Environment) Push() *Environment {
	return &Environment{
		Parent:    env,
		HasParent: true,
		Env:       make(map[string]*types.Type),
	}
}

func (env *Environment) Set(key string, value *types.Type) {
	env.Env[key] = value
}

func (env *Environment) Get(str string) (*types.Type, bool) {
	return env.get(str, 0)
}

func (env *Environment) get(str string, r int) (t *types.Type, ok bool) {

	if _, ok := env.Env[str]; ok {
		return env.Env[str], true
	}

	if !env.HasParent {
		return t, false
	}

	return env.Parent.get(str, r+1)
}
