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

	//ast.Print(fset, obj)

	//switch x := obj.Decl.(type) {
	//case *ast.TypeSpec:
	//	switch y := x.Type.(type) {
	//	case *ast.InterfaceType:
	//		//ast.Print(fset, y)
	//
	//		for _, x := range y.Methods.List {
	//			ast.Print(fset, x)
	//		}
	//	}
	//}

	//// Inspect the AST and print all identifiers and literals.

	mock.Methods = make(map[string]Method)

	var currentMethodName string

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.InterfaceType:
			fmt.Println("hit an interface")
			//ast.Print(fset, x)

			for _, method := range x.Methods.List {
				fmt.Println(method)

				for _, x := range method.Names {
					fmt.Println("method name", x)
					currentMethodName = x.Name

					ast.Inspect(method, func(n ast.Node) bool {
						switch x := n.(type) {
						case *ast.FieldList:
							for _, field := range x.List {
								field.Type
								fmt.Println("i want to add", field, "as arguments to", currentMethodName)

							}
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
