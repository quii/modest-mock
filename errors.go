package modest_mock

import "fmt"

type InterfaceNotFoundError struct {
	Name string
}

func (i *InterfaceNotFoundError) Error() string {
	return fmt.Sprintf("Unable to find interface %s in source", i.Name)
}

type ParseFailError struct {
	err error
}

func (p *ParseFailError) Error() string {
	return fmt.Sprintf("Unable to parse source code %s", p.err.Error())
}
