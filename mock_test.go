package modest_mock

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	scenarios := []struct {
		Name         string
		Src          string
		ExpectedMock Mock
		ExpectedErr  error
	}{
		{
			Name: "Interface with named arguments and return values",
			Src: `
				package main
				type Store interface{
					Save(firstname, lastname string) (err error)
				}
`,
			ExpectedMock: Mock{
				Name: "Store",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"firstname", "string"},
							{"lastname", "string"},
						},
						ReturnValues: []Value{
							{"err", "error"},
						},
					},
				},
			},
			ExpectedErr: nil,
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			mock, err := New(strings.NewReader(s.Src))

			assert.Equal(t, s.ExpectedErr, err)
			assert.Equal(t, s.ExpectedMock, mock)
		})
	}
}
