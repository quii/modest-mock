package modestmock

import (
	"bytes"
	"text/template"
)

const mockTemplate = `package {{.Package}}

type {{.Name}}Mock struct {
}
`

func GenerateMockCode(mock Mock) (string, error) {
	tmpl, err := template.New("test").Parse(mockTemplate)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, mock)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
