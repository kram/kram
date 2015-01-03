package main

type String struct {
	String bool
	value string
}

func (s *String) Init(str string) bool {
	s.value = str

	return true
}

func (s *String) Add(str string) {
	s.value = s.value + str
}

func (s *String) toString() string {
	return s.value
}