package modest_mock

import "io"

type Value struct {
	Name, Type string
}

type Method struct {
	Arguments    []Value
	ReturnValues []Value
}

type Mock struct {
	Name    string
	Methods map[string]Method
}

func New(src io.Reader) (Mock, error) {
	return Mock{}, nil
}
