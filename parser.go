package modestmock

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strconv"
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

	foundInterface := false

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {

		case *ast.Ident:
			foundInterface = x.Name == name

		case *ast.InterfaceType:
			if foundInterface {
				for _, method := range x.Methods.List {
					addMethod(method, method.Names[0].Name, mock, fset)

				}
				foundInterface = false
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

func addMethod(method *ast.Field, name string, mock Mock, fset *token.FileSet) {
	ast.Inspect(method, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncType:

			m := Method{
				Arguments: getValues(x.Params, fset),
			}

			if x.Results != nil {
				m.ReturnValues = getValues(x.Results, fset)
			}

			mock.Methods[name] = m
		}
		return true
	})
}

func getValues(list *ast.FieldList, fset *token.FileSet) (values []Value) {
	for _, field := range list.List {

		fieldType := getType(field, fset)

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

func getType(field *ast.Field, fset *token.FileSet) string {

	if arr, isArr := field.Type.(*ast.ArrayType); isArr {

		len := 0

		if size, isFixedLen := arr.Len.(*ast.BasicLit); isFixedLen {
			len, _ = strconv.Atoi(size.Value)
		}

		arrDelc := "[]"

		if len > 0 {
			arrDelc = fmt.Sprintf("[%d]", len)
		}

		if basic, isBasic := arr.Elt.(*ast.Ident); isBasic {
			return fmt.Sprintf("%s%s", arrDelc, basic.Name)
		}

		if complex, isComplex := arr.Elt.(*ast.SelectorExpr); isComplex {
			pkg := complex.X.(*ast.Ident)
			typ := complex.Sel.Name
			return fmt.Sprintf("%s%s.%s", arrDelc, pkg.Name, typ)
		}

	}

	if complex, isComplexType := field.Type.(*ast.SelectorExpr); isComplexType {
		pkg := complex.X.(*ast.Ident)
		typ := complex.Sel.Name
		return fmt.Sprintf("%s.%s", pkg.Name, typ)
	}
	return fmt.Sprintf("%v", field.Type)
}
