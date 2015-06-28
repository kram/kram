// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package environment

import (
	"github.com/kram/kram/src/types"
)

type Environment struct {
	Env       map[string]*types.Value
	HasParent bool
	Parent    *Environment
}

func (env *Environment) Init() {
	env.Env = make(map[string]*types.Value)
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
		Env:       make(map[string]*types.Value),
	}
}

func (env *Environment) Set(key string, value *types.Value) {
	env.Env[key] = value
}

func (env *Environment) Get(str string) (*types.Value, bool) {
	return env.get(str, 0)
}

func (env *Environment) get(str string, r int) (t *types.Value, ok bool) {

	if _, ok := env.Env[str]; ok {
		return env.Env[str], true
	}

	if !env.HasParent {
		return t, false
	}

	return env.Parent.get(str, r+1)
}
