package types

type Type interface {
	Init(string)
	Math(string, Type) Type
	Compare(string, Type) Type
	Type() string
	ToString() string
}
