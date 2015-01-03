package main

type Bool struct {
	Bool bool
	value bool
}

func (b *Bool) Init(str string) bool {
	
	if str == "true" {
		b.value = true
	} else {
		b.value = false
	}

	return true
}

func (b *Bool) toString() string {

	if b.value {
		return "true"
	}

	return "false"
}