package modestmock

import "fmt"

type Value struct {
	Name, Type string
}

func (v Value) AsCodeDeclaration() string {
	return fmt.Sprintf("%s %s", v.Name, v.Type)
}

type Method struct {
	Arguments    []Value
	ReturnValues []Value
}

type Mock struct {
	Name    string
	Imports []string
	Package string
	Methods map[string]Method
}

//todo: test me
func (m *Mock) ReturnValues() map[string][]Value {
	returns := make(map[string][]Value)

	for method, values := range m.Methods {
		if len(values.ReturnValues) > 0 {
			returns[method] = values.ReturnValues
		}
	}

	return returns
}

func (m *Mock) HasReturnValues() bool {
	return len(m.ReturnValues()) > 0
}

func (m *Mock) Arguments() map[string][]Value {
	arguments := make(map[string][]Value)

	for method, values := range m.Methods {
		arguments[method] = values.Arguments
	}

	return arguments
}
