# modest mock

modest mock is a tool to create mocks from go interfaces for your tests

# Aims

- *Simple* - the usage is just `PATH_TO_INTERFACE` `NAME_OF_INTERFACE`
- *Type safe* - You are using Go because you like static checking, you
  shouldn't have to give that up in your tests so no reflection or
`interface{}`
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

if mockStore.Calls.Save.firstname != "expected firstname" {
  t.Error("Didnt't call Store with correct firstname")
}
```

- Easy to set return values too

`mockStore.Returns.Save.err = errors.New("Simulating save failure")`

I haven't written anything yet, but aspirations are nice.
