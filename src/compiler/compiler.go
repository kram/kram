package compiler

import (
	"bytes"
	"log"
	"os"

	ins "github.com/kram/kram/src/instructions"
	"github.com/kram/kram/src/types"
)

type compiler struct {
	output *bytes.Buffer
}

func Run(tree ins.Block) {
	c := &compiler{}

	c.output = bytes.NewBufferString("")

	c.bootstrapStart()

	c.Operation(tree, types.ON_NOTHING)

	c.bootstrapEnd()

	c.output.WriteTo(os.Stdout)
}

func (c *compiler) Operation(node ins.Node, on types.ON) {

	if push, ok := node.(ins.PushClass); ok {
		c.pushClass(push)
		return
	}

	if call, ok := node.(ins.Call); ok {
		c.call(call)
		return
	}

	if literal, ok := node.(ins.Literal); ok {
		c.literal(literal)
		return
	}

	if variable, ok := node.(ins.Variable); ok {
		c.variable(variable)
		return
	}

	if set, ok := node.(ins.Set); ok {
		c.set(set)
		return
	}

	if def, ok := node.(ins.DefineFunction); ok {
		c.defineFunction(def)
		return
	}

	if call, ok := node.(ins.Call); ok {
		c.call(call)
		return
	}

	if block, ok := node.(ins.Block); ok {
		for _, body := range block.Body {
			c.Operation(body, types.ON_NOTHING)
		}

		return
	}

	log.Panic("Can not handle type", node)
}

func (c *compiler) operationList(nodes []ins.Node) {
	for _, node := range nodes {
		c.Operation(node, types.ON_NOTHING)
	}
}

func (c *compiler) pushClass(i ins.PushClass) {
	c.Operation(i.Left, types.ON_NOTHING)
	c.output.WriteString(".")
	c.Operation(i.Right, types.ON_NOTHING)
}

func (c *compiler) call(i ins.Call) {
	c.Operation(i.Left, types.ON_NOTHING)

	c.output.WriteString("(")

	for _, arg := range i.Arguments {
		c.Operation(arg.Value, types.ON_NOTHING)
	}

	c.output.WriteString(")")
}

func (c *compiler) literal(lit ins.Literal) {
	if lit.Type == "string" {
		c.output.WriteString("\"" + lit.Value + "\"")
		return
	}

	if lit.Type == "number" {
		c.output.WriteString(lit.Value)
		return
	}

	log.Panic("Unknown literal type: ", lit.Type)
}

func (c *compiler) variable(v ins.Variable) {
	c.output.WriteString(v.Name)
}

func (c *compiler) set(v ins.Set) {
	c.output.WriteString(v.Name + " := ")
	c.Operation(v.Right, types.ON_NOTHING)
}

func (c *compiler) defineFunction(def ins.DefineFunction) {
	c.output.WriteString("func() {\n")
	c.operationList(def.Body.Body)
	c.output.WriteString("}\n\n")
}

func (c *compiler) bootstrapStart() {
	c.output.WriteString(`package main
    
import (
    IO "fmt"
)

func main() {
`)
}

func (c *compiler) bootstrapEnd() {
	c.output.WriteString("\n}")
}
