package modestmock

import (
	"bytes"
	"strings"
	"text/template"
)

const mockStructTemplate = `package {{.Package}}

type {{.Name}}Mock struct {
{{ if .HasReturnValues -}}
	Returns struct {
		{{ range $method, $arguments := .ReturnValues -}}
		{{ $method }} struct {
			{{range $arg := $arguments -}}
			{{ $arg.Name }} {{ $arg.Type }}
			{{end -}}
		}
		{{end -}}
	}
{{end -}}
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
{{.Returns}}
}
`

func generateMethod(receiver string, methodName string, method Method) (string, error) {
	tmpl, err := template.New("struct").Parse(mockMethodTemplate)

	if err != nil {
		return "", err
	}

	receiverVarName := strings.ToLower(string(receiver[0]))

	returnStatement := ""
	if len(method.ReturnValues) > 0 {
		var returns []string
		for _, r := range method.ReturnValues {
			returns = append(returns, receiverVarName+".Returns."+methodName+"."+r.Name)
		}
		returnStatement = "\treturn " + strings.Join(returns, ", ")
	}

	viewModel := struct {
		ReceiverVar, Receiver, Name, Arguments, Returns string
	}{
		ReceiverVar: receiverVarName,
		Receiver:    receiver,
		Name:        methodName,
		Arguments:   generateFields(method.Arguments),
		Returns:     returnStatement,
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, viewModel)

	return buffer.String(), nil

}
