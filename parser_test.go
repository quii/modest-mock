package modestmock

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
			Name:          "No return values, one arg",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(age int)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"age", "int"},
						},
					},
				},
			},
		},

		{
			Name:          "Slices",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(numbers []int)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"numbers", "[]int"},
						},
					},
				},
			},
		},

		{
			Name:          "Foreign Slices",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(numbers []special.Number)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"numbers", "[]special.Number"},
						},
					},
				},
			},
		},

		{
			Name:          "Foreign arrays",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(numbers [2]special.Number)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"numbers", "[2]special.Number"},
						},
					},
				},
			},
		},

		{
			Name:          "Extraneous interfaces",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(age int)
						}

						type NotImportant interface{
							Yell()
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"age", "int"},
						},
					},
				},
			},
		},

		{
			Name:          "No return values, two args same type",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(age, height int)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"age", "int"},
							{"height", "int"},
						},
					},
				},
			},
		},

		{
			Name:          "No return values, two args different types",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(age int, name string)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"age", "int"},
							{"name", "string"},
						},
					},
				},
			},
		},

		{
			Name:          "Named return values",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Delete(id int) (success bool)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Delete": {
						Arguments: []Value{
							{"id", "int"},
						},
						ReturnValues: []Value{
							{"success", "bool"},
						},
					},
				},
			},
		},

		{
			Name:          "Anonymous return values",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Delete(id int) bool
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Delete": {
						Arguments: []Value{
							{"id", "int"},
						},
						ReturnValues: []Value{
							{"", "bool"},
						},
					},
				},
			},
		},
		{
			Name:          "Non primative types with imports",
			InterfaceName: "Channel",
			Src: `
						package main

						import "github.com/streadway/amqp"

						type Channel interface{
							QueueDeclare(name string, args amqp.Table) (amqp.Queue, error)
						}
		`,
			ExpectedMock: Mock{
				Name:    "Channel",
				Package: "main",
				Imports: []string{`"github.com/streadway/amqp"`},
				Methods: map[string]Method{
					"QueueDeclare": {
						Arguments: []Value{
							{"name", "string"},
							{"args", "amqp.Table"},
						},
						ReturnValues: []Value{
							{"", "amqp.Queue"},
							{"", "error"},
						},
					},
				},
			},
		},

		{
			Name:          "Multiple methods",
			InterfaceName: "Store",
			Src: `
						package main
						type Store interface{
							Save(firstname, lastname string)
							Delete(id int) bool
						}
		`,
			ExpectedMock: Mock{
				Name:    "Store",
				Package: "main",
				Methods: map[string]Method{
					"Save": {
						Arguments: []Value{
							{"firstname", "string"},
							{"lastname", "string"},
						},
					},
					"Delete": {
						Arguments: []Value{
							{"id", "int"},
						},
						ReturnValues: []Value{
							{"", "bool"},
						},
					},
				},
			},
		},

		{
			Name:        "Invalid go code returns an error",
			Src:         `function poo() { console.log("lolz"); }`,
			ExpectError: &ParseFailError{},
		},

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
