package main

type Lib interface {
	Init(*VM, []Type)
	Setup()
	Instance() (Lib, string)
	ToString() string
}

type Library struct{}

func (lib *Library) Init(vm *VM, params []Type) {}
func (lib *Library) Setup()                     {}

func (lib *Library) ToString() string {
	return "DEFAULT_LIBRARY"
}

func (lib *Library) Instance() (Lib, string) {
	return &Library{}, "DEFAULT_LIBRARY"
}

func DefaultReturn() Type {
	bl := Bool{}
	bl.Init("false")

	return &bl
}
