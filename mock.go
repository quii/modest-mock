package modestmock

type Value struct {
	Name, Type string
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
