// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package kram

import (
	ins "github.com/kram/kram/src/instructions"
)

// This is the stack for the parser
// It holds the current state of the parser, and is used to know what we're currently parsing
type Stack struct {
	Items   *[]ins.Node
	Parents []*[]ins.Node
}

func (stack *Stack) Pop() {
	if len(stack.Parents) == 0 {
		items := make([]ins.Node, 0)
		stack.Items = &items
		return
	}

	stack.Items = stack.Parents[len(stack.Parents)-1]
	stack.Parents = stack.Parents[:len(stack.Parents)-1]
}

func (stack *Stack) Push() {
	stack.Parents = append(stack.Parents, stack.Items)

	items := make([]ins.Node, 0)
	stack.Items = &items
}

func (stack *Stack) Add(node ins.Node) {
	items := *stack.Items
	items = append(items, node)

	stack.Items = &items
}

func (stack *Stack) Reset() {
	stack.Empty()
	stack.Parents = make([]*[]ins.Node, 0)
}

func (stack *Stack) Empty() {
	items := make([]ins.Node, 0)
	stack.Items = &items
}
