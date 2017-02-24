package modestmock

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

func New(src io.Reader, name string) (mock Mock, err error) {

	fset := token.NewFileSet()
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
	mock.Imports = getImports(f.Imports)
	mock.Package = f.Name.Name

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.InterfaceType:
			for _, method := range x.Methods.List {
				addMethod(method, method.Names[0].Name, mock)
			}
		}
		return true
	})

	return
}

func getImports(importSpec []*ast.ImportSpec) (imports []string) {
	for _, i := range importSpec {
		imports = append(imports, i.Path.Value)
	}
	return
}

func addMethod(method *ast.Field, name string, mock Mock) {
	ast.Inspect(method, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncType:

			m := Method{
				Arguments: getValues(x.Params),
			}

			if x.Results != nil {
				m.ReturnValues = getValues(x.Results)
			}

			mock.Methods[name] = m
		}
		return true
	})
}

func getValues(list *ast.FieldList) (values []Value) {
	for _, field := range list.List {

		fieldType := getType(field)

		if len(field.Names) == 0 {
			values = append(values, Value{
				"", fieldType,
			})
		}

		for _, f := range field.Names {
			values = append(values, Value{
				f.Name, fieldType,
			})
		}
	}
	return
}

func getType(field *ast.Field) string {
	if complex, isComplexType := field.Type.(*ast.SelectorExpr); isComplexType {
		pkg := complex.X.(*ast.Ident)
		typ := complex.Sel.Name
		return fmt.Sprintf("%s.%s", pkg.Name, typ)
	}
	return fmt.Sprintf("%v", field.Type)
}
