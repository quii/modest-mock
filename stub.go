package modest_mock

type Value struct {
	Name, Type string
}

type Method struct {
	Arguments    []Value
	ReturnValues []Value
}

type Stub struct {
	Name    string
	Methods map[string]Method
}
