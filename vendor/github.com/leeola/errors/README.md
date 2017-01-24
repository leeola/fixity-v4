
# errors (unstable & incomplete)

leeola/errors is a simple package meant to summarize an errors history as it
propagates through your program. As you wrap errors, it will keep track of each
file location and the added error context. The package has various methods to extract
this information in meaningful ways. `errors.Println(err)` and
`s := errors.Sprintln(err)` for example.

If the originating error *(the one not from leeola/errors)* is needed, `errors.Cause(err)`
can be used to extract the original error.

## Example

    ~/g/s/g/l/errors> go run _example/example.go
    Baz() returned an error.

    err.Error():
    Bar() failed: Foo() failed: os.Open failed: open fakefile.txt: no such file or directory

    errors.Println(err):
    errors/_example/example.go:13: os.Open failed: open fakefile.txt: no such file or directory
    errors/_example/example.go:21: Foo() failed
    errors/_example/example.go:29: Bar() failed

    errors.Cause(err).Error():
    open fakefile.txt: no such file or directory

