<div align="center">

# go-aloe
[![Continuous Integration](https://img.shields.io/github/actions/workflow/status/tfadeyi/go-aloe/ci.yml?branch=main&style=flat-square)](https://github.com/tfadeyi/go-aloe/actions/workflows/ci.yml)

</div>

---

Aloe is an easy-to-use Go library for wrapping and handling errors. The library allows applications to define errors
in an [Aloe specification](schema/README.md), which is stored at the application root level.
This can then be wrapped to existing error present in an application source code, giving additional information to
the application users.

> ⚠ Currently, under development.

> ⚠ The project is currently missing the aloe-cli for easier generation of the specification and the static website.  

## Features

- Wraps the application errors with errors defined in the aloe specification.
- Link errors to specific error pages in on static website. (needs aloe-cli to be released) 
- Supports Aloe specification in yaml,json and toml formats.

## Installation

To use **go-aloe**, you need to have Go installed. Then, you can install the library using the following command:

```bash
go get -u github.com/tfadeyi/go-aloe
```

## Usage
Define an Aloe specification for your application:
>**Save the specification, under the application current working directory, with `default.aloe.toml` as the file name.**

```toml
BaseUrl = "https://github.com/tfadeyi/my-app"
Description = "Sample application"
Name = "my-app"
Title = "My Application"
Version = "v0.0.1"

[ErrorsDefinitions]
    [ErrorsDefinitions.error_something_code]
    code = "error_something_code"
    summary = "This is a summary of the error that will wrap the application error."
    title = "Error On Something"
```

For more info on the Aloe specification check the schema [README](schema/README.md).

Import the library into your Go code:

```go
import "github.com/tfadeyi/go-aloe"
```

Now, you can use the library's functions and types in your code. Here's a simple example:

```go
    package main
    
    import (
        "errors"
        "log"
    
        "github.com/tfadeyi/go-aloe"
    )
    
    func doSomething() error {
        err := errors.New("something")
        return goaloe.DefaultOrDie().Error(err, "error_something_code")
    }
    func main() {
        err := doSomething()
        if err != nil {
            log.Fatal(err)
        }
    }
```

Running the code will result:

```text
$ go run main.go
This is a summary of the error that will wrap the application error.
for additional info check https://github.com/tfadeyi/my-app/my-app/error_something_code: [something]
```

## License
Go-aloe is released under the [MIT](./LICENSE) License.
