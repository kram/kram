package gus

func DefaultReturn() Type {
	bl := Bool{}
	bl.Init("false")

	return &bl
}
