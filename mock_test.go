package modest_mock

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	scenarios := []struct {
		Name          string
		InterfaceName string
		Src           string
		ExpectedMock  Mock
		ExpectError   error
	}{
		{
			Name:        "Invalid go code returns an error",
			Src:         `function poo() { console.log("lolz"); }`,
			ExpectError: &ParseFailError{},
		},
		//{
		//	Name:          "Valid go code but interface missing",
		//	InterfaceName: "NotStore",
		//	Src: `
		//				package main
		//				type Store interface{
		//					Save(firstname, lastname string) (err error)
		//				}
		//`,
		//	ExpectedMock: Mock{
		//		Name: "Store",
		//		Methods: map[string]Method{
		//			"Save": {
		//				Arguments: []Value{
		//					{"firstname", "string"},
		//					{"lastname", "string"},
		//				},
		//				ReturnValues: []Value{
		//					{"err", "error"},
		//				},
		//			},
		//		},
		//	},
		//},
		{
			Name:          "Valid go code but interface missing",
			InterfaceName: "NotStore",
			Src: `
						package main
						type Store interface{
							Save(firstname, lastname string) (err error)
						}
		`,
			ExpectError: &InterfaceNotFoundError{},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			mock, err := New(strings.NewReader(s.Src), s.InterfaceName)

			if s.ExpectError != nil {
				assert.IsType(t, s.ExpectError, err)
			}

			assert.Equal(t, s.ExpectedMock, mock)
		})
	}
}
