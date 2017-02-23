# tips from Fatih Arslan

- https://speakerdeck.com/farslan/how-to-write-your-own-go-tool it was also recorded and you can find it here: https://youtu.be/oxc8B2fjDvY

- https://medium.com/@farslan/a-look-at-go-scanner-packages-11710c2655fc#.73dufutkh

- https://github.com/fatih/gomodifytags (project that uses AST a lot)

- If you want to generate code, look also at the `go/types` package. Here is a farly detailed, but still a complex tutorial: https://github.com/golang/example/tree/master/gotypes

# general notes

- using go generate seems like a good plan here, people can put the directive on top of the interface they want to generate mocks on
- https://github.com/golang/example/tree/master/gotypes use this to check the validity of the code generated

## open questions

- need to do imports if we're going to validate code, but how do i know which ones aren't needed?