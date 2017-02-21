package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// https://blog.gopheracademy.com/code-generation-from-the-ast/

func main() {
	src := `
package main
type Store interface{
	Save(firstname, lastname string) (err error)
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
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
}
