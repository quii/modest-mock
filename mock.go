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

	if obj==nil{
		return mock, &InterfaceNotFoundError{name}
	}

	fmt.Println("blah")
	ast.Print(fset, obj)

	switch x:= obj.Decl.(type){
	case *ast.TypeSpec:
		switch y := x.Type.(type) {
		case *ast.InterfaceType:
			ast.Print(fset, y)
		}
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.InterfaceType:
			ast.Print(fset, x)

			//for _, method := range x.Methods.List {
			//	fmt.Println(method)
			//
			//	for _, x := range method.Names {
			//		fmt.Println(x)
			//	}
			//}
			//case *ast.FieldList:
			//	fmt.Println("number of fields", x.NumFields())
			//	for _, field := range x.List {
			//		fmt.Println(field)
			//	}
		}
		return true
	})

	return
}
