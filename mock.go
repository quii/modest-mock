package modest_mock

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

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

func New(src io.Reader, name string) (mock Mock, err error) {

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return mock, &ParseFailError{err}
	}

	obj := f.Scope.Lookup(name)

	if obj == nil {
		return mock, &InterfaceNotFoundError{name}
	}

	mock.Methods = make(map[string]Method)
	mock.Name = name

	var currentMethodName string

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.InterfaceType:
			for _, method := range x.Methods.List {

				for _, x := range method.Names {
					currentMethodName = x.Name

					ast.Inspect(method, func(n ast.Node) bool {
						switch x := n.(type) {
						case *ast.FuncType:

							m := Method{
								Arguments: getValues(x.Params),
							}

							if x.Results != nil {
								m.ReturnValues = getValues(x.Results)
							}

							mock.Methods[currentMethodName] = m
						}
						return true
					})
				}

			}
			return true

		}
		return true
	})

	return
}

func getValues(list *ast.FieldList) (values []Value) {
	for _, field := range list.List {
		for _, f := range field.Names {
			values = append(values, Value{
				f.Name, fmt.Sprintf("%v", field.Type),
			})
		}
	}
	return
}
