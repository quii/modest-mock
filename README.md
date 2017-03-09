# modest mock

modest mock is a tool to generate type-safe mocks from go interfaces.

# Aims

- *Simple* - the usage is just `PATH_TO_INTERFACE` `NAME_OF_INTERFACE`
- *Type safe* - You are using Go because you like static checking, you shouldn't have to give that up in your tests so no reflection, magic strings or `interface{}`
- *Not a framework*, it does not prescribe how you write your tests
- No dependencies required in your app, just generate the code once.
- Easy to use assertions

Given

```go
type Store interface{
  Save(firstname, lastname string) (err error)
}
```

I would expect to use the mock in my test like:

```go
mockStore := NewMockStore()

thingUnderTest := NewThing(mockStore)
thing.DoIt("expected firstname", "smith")

if mockStore.Calls.Save[0].firstname != "expected firstname" {
  t.Error("Didnt't call Store with correct firstname on the first call")
}
```

## Usage

`go get github.com/quii/modest-mock/cmd/modestmock`

### On the command line

`modestmock -path=$PATH_TO_FILE_WITH_INTERFACE -name=$NAME_OF_INTERFACE`

This will send the generated code to stdout so you can pipe it to a file.

### Using go generate

`//go:generate modestmock -name=Bank -out=bank_mock.go`

Have a look at the `example` directory for a real world-ish example. 