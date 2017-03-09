package modestmock

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

const mockStructTemplate = `
{{ $mock := . }}

package {{.Package}}

import "fmt"

type {{.Name}}Mock struct {

Calls struct {
{{ range $name, $method := .Methods }}
	{{ $name }} []{{$mock.Name}}Mock_{{$name}}Args
{{end}}
}

{{ if .HasReturnValues }}
Returns struct {
	{{ range $name, $method := .Methods }}
		{{ $name }} map[{{$mock.Name}}Mock_{{$name}}Args]{{$mock.Name}}Mock_{{$name}}Returns
	{{end}}
}
{{end}}
}
`

func GenerateMockCode(mock Mock) (string, error) {
	receiver := mock.Name + "Mock"

	mockStruct, err := generateMockStruct(mock)

	if err != nil {
		return "", err
	}

	var methods []string
	for name, definition := range mock.Methods {
		stubMethod, err := generateMethod(receiver, name, definition)

		if err != nil {
			return "", err
		}

		methods = append(methods, stubMethod)
	}

	allMethods := strings.Join(methods, "\n")

	methodTypes, err := generateArgAndReturnTypes(mock)

	if err != nil {
		return "", err
	}

	constructor, err := generateConstructor(mock)

	if err != nil {
		return "", err
	}

	code := mockStruct + constructor + allMethods + methodTypes

	formattedCode, err := format.Source([]byte(code))

	return string(formattedCode), err

}

const constructorTemplate = `
{{ $mock := . }}
func New{{.Name}}Mock() *{{.Name}}Mock {
	newMock := new({{.Name}}Mock)

	{{ range $name, $method := .Methods }}{{/*
		*/}}newMock.Returns.{{ $name }} = make(map[{{$mock.Name}}Mock_{{$name}}Args]{{$mock.Name}}Mock_{{$name}}Returns)
	{{end}}
	return newMock
}
`

func generateConstructor(mock Mock) (string, error) {
	tmpl, err := template.New("constructor").Parse(constructorTemplate)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	err = tmpl.Execute(&buffer, &mock)

	if err != nil {
		return "", err
	}

	return buffer.String(), err
}

const methodTypeTemplate = `
type {{.Name}}Mock_{{.MethodName}}Args struct {
		{{range $arg := .Method.Arguments }}{{/*
			*/}}{{- $arg.AsCodeDeclaration }}
		{{end}}
}

type {{.Name}}Mock_{{.MethodName}}Returns struct {
		{{range $arg := .Method.ReturnValues }}{{/*
			*/}}{{ $arg.AsCodeDeclaration }}
		{{end}}
}

`

func generateArgAndReturnTypes(mock Mock) (string, error) {
	tmpl, err := template.New("methodType").Parse(methodTypeTemplate)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	for name, method := range mock.Methods {
		viewModel := struct {
			Name       string
			MethodName string
			Method     Method
		}{
			mock.Name,
			name,
			method,
		}
		err = tmpl.Execute(&buffer, viewModel)

		if err != nil {
			return "", err
		}
	}

	return buffer.String(), nil
}

func generateMockStruct(mock Mock) (string, error) {
	tmpl, err := template.New("struct").Parse(mockStructTemplate)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &mock)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func generateFields(values []Value) string {
	var args []string

	for _, v := range values {
		args = append(args, v.Name+" "+v.Type)
	}

	return "(" + strings.Join(args, ",") + ")"
}

const mockMethodTemplate = `
func ({{.ReceiverVar}} *{{.Receiver}}) {{.Name}}{{.Arguments}} {{.ReturnArgs}} {
{{.RecordCall}}
{{.ReturnStatement}}
}
`

func generateMethod(receiver string, methodName string, method Method) (string, error) {
	tmpl, err := template.New("struct").Parse(mockMethodTemplate)

	if err != nil {
		return "", err
	}

	receiverVarName := strings.ToLower(string(receiver[0]))

	returnStatement := ""
	returnArgs := ""

	if len(method.ReturnValues) > 0 {
		returnStatement, err = generateReturn(receiverVarName, methodName, method.ReturnValues)

		if err != nil {
			return "", err
		}

		returnArgs = generateFields(method.ReturnValues)
	}

	viewModel := struct {
		ReceiverVar, Receiver, Name, Arguments, ReturnStatement, ReturnArgs, RecordCall string
	}{
		ReceiverVar:     receiverVarName,
		Receiver:        receiver,
		Name:            methodName,
		Arguments:       generateFields(method.Arguments),
		ReturnStatement: returnStatement,
		ReturnArgs:      returnArgs,
		RecordCall:      generateRecordCall(receiverVarName, receiver, methodName, method.Arguments),
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, viewModel)

	return buffer.String(), nil

}

const returnStatementTmp = `
	if vals, exists := {{.ReceiverName}}.Returns.{{.MethodName}}[call]; exists {
		return {{.ReturnStmt}}
	}

	panic(fmt.Sprintf("no return values found for args %+v, ive got %+v", call, {{.ReceiverName}}.Returns.{{.MethodName}}))`

func generateReturn(receiverName, methodName string, returnVals []Value) (string, error) {
	tmpl, err := template.New("returns").Parse(returnStatementTmp)

	if err != nil {
		return "", err
	}

	var returnValsWithRecievers []string

	for _, v := range returnVals {
		returnValsWithRecievers = append(returnValsWithRecievers, fmt.Sprintf("vals.%s", v.Name))
	}

	var buffer bytes.Buffer

	viewModel := struct {
		ReceiverName string
		MethodName   string
		ReturnStmt   string
	}{
		receiverName,
		methodName,
		strings.Join(returnValsWithRecievers, ","),
	}
	err = tmpl.Execute(&buffer, &viewModel)

	return buffer.String(), err
}

func generateRecordCall(recieverVar, reciever, method string, arguments []Value) string {
	var allFields []string
	var allValues []string
	for _, arg := range arguments {
		allValues = append(allValues, arg.Name)
		allFields = append(allFields, arg.AsCodeDeclaration())
	}

	// {{.Name}}Mock_{{.MethodName}}Args
	methodArgsType := reciever + "_" + method + "Args"

	call := fmt.Sprintf("call := %s{%s}", methodArgsType, strings.Join(allValues, ","))
	callToUpdate := fmt.Sprintf("%s.Calls.%s", recieverVar, method)
	updateCall := fmt.Sprintf("%s = append(%s, call)", callToUpdate, callToUpdate)

	code := fmt.Sprintf("%s\n%s", call, updateCall)

	return code
}
