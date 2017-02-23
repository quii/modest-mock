package modest_mock

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"errors"
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

var badCodeErr = errors.New("Bad code you stupid poo")

func New(src io.Reader, name string) (mock Mock, err error) {

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		return
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.InterfaceType:
			fmt.Println(x.Interface)
			for _, method := range x.Methods.List {
				fmt.Println(method)

				for _, x := range method.Names {
					fmt.Println(x)
				}
			}
		case *ast.FieldList:
			fmt.Println("number of fields", x.NumFields())
			for _, field := range x.List {
				fmt.Println(field)
			}
		}
		return true
	})

	return
}
