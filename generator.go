package modestmock

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
)

const mockStructTemplate = `

package {{.Package}}

type {{.Name}}Mock struct {

Calls struct {
{{ range $name, $method := .Methods }}
	{{ $name }} []struct {
		{{range $arg := $method.Arguments }}
			{{ $arg.AsCodeDeclaration }}
		{{end}}
	}
{{end}}
}

{{ if .HasReturnValues }}
Returns struct {
	{{ range $name, $method := .Methods }}
		{{ $name }} []struct {
			{{range $arg := $method.ReturnValues }}
				{{ $arg.AsCodeDeclaration }}
			{{end}}
		}
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

	code := mockStruct + allMethods

	formattedCode, err := format.Source([]byte(code))

	return string(formattedCode), err
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
		var returns []string
		for _, r := range method.ReturnValues {
			returnIndex := "[len(" + receiverVarName + ".Calls." + methodName + ")-1]"
			returns = append(returns, receiverVarName+".Returns."+methodName+returnIndex+"."+r.Name)
		}
		returnStatement = "\treturn " + strings.Join(returns, ", ")

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
		RecordCall:      generateRecordCall(receiverVarName, methodName, method.Arguments),
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, viewModel)

	return buffer.String(), nil

}

func generateRecordCall(reciever string, method string, arguments []Value) string {
	var allFields []string
	var allValues []string
	for _, arg := range arguments {
		allValues = append(allValues, arg.Name)
		allFields = append(allFields, arg.AsCodeDeclaration())
	}

	call := fmt.Sprintf("call := struct{%s}{%s}", strings.Join(allFields, "\n"), strings.Join(allValues, ","))
	callToUpdate := fmt.Sprintf("%s.Calls.%s", reciever, method)
	updateCall := fmt.Sprintf("%s = append(%s, call)", callToUpdate, callToUpdate)

	code := fmt.Sprintf("%s\n%s", call, updateCall)

	return code
}
