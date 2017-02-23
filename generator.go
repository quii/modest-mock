package modestmock

import (
	"bytes"
	"strings"
	"text/template"
)

const mockStructTemplate = `package {{.Package}}

type {{.Name}}Mock struct {
{{- if .HasReturnValues }}
	Returns struct {
		ffs int
	}
{{end}}
}
`

func GenerateMockCode(mock Mock) (string, error) {
	receiver := mock.Name + "Mock"

	mockStruct, err := generateMockStruct(mock)

	var methods []string
	for name, definition := range mock.Methods {
		stubMethod, err := generateMethod(receiver, name, definition)

		if err != nil {
			return "", err
		}

		methods = append(methods, stubMethod)
	}

	allMethods := strings.Join(methods, "\n")

	return mockStruct + allMethods, err
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
func ({{.ReceiverVar}} *{{.Receiver}}) {{.Name}}{{.Arguments}} {

}
`

func generateMethod(receiver string, methodName string, method Method) (string, error) {
	tmpl, err := template.New("struct").Parse(mockMethodTemplate)

	if err != nil {
		return "", err
	}

	viewModel := struct {
		ReceiverVar, Receiver, Name, Arguments string
	}{
		ReceiverVar: strings.ToLower(string(receiver[0])),
		Receiver:    receiver,
		Name:        methodName,
		Arguments:   generateFields(method.Arguments),
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, viewModel)

	return buffer.String(), nil

}
