<!-- markdownlint-configure-file {
  "MD013": {
    "code_blocks": false,
    "tables": false
  },
  "MD033": false,
  "MD041": false
} -->

<div align="center">

<img title="error.fyi" alt="error.fyi" src="./docs/logo.png">

<p>
    <strong>[<a href="https://docs.error.fyi">Docs</a>]</strong>
    <strong>[<a href="#-features">Features</a>]</strong>
    <strong>[<a href="#-installation">Installation</a>]</strong>
    <strong>[<a href="#-get-started">Get Started</a>]</strong>
</p>

</div>

#### Generate and maintain user facing error documentation from in-code annotations.

---
**error.fyi** is a library and CLI tool for Golang applications. It lets you generate targeted error documentation for user facing errors
which can be made available through a static site generator and hosting services.
The **error.fyi** library will then wrap application errors with additional context and link to the target error page.

[![GitHub release](https://img.shields.io/github/v/release/tfadeyi/errors?color=green&style=for-the-badge)](https://github.com/tfadeyi/errors/releases)
[![Nix Devshell](https://img.shields.io/badge/nix-devshell-blue?logo=NixOS&style=for-the-badge)](https://github.com/tfadeyi/errors)
[![Continuous Integration](https://img.shields.io/github/actions/workflow/status/tfadeyi/errors/ci.yml?branch=main&style=for-the-badge)](https://github.com/tfadeyi/errors/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-yellowgreen.svg?style=for-the-badge)](https://github.com/tfadeyi/errors/blob/main/LICENSE)
[![Language](https://img.shields.io/github/go-mod/go-version/tfadeyi/errors?style=for-the-badge)](https://github.com/tfadeyi/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/tfadeyi/errors?style=for-the-badge)](https://goreportcard.com/report/github.com/tfadeyi/errors)

> **Note**
> The project is not yet production ready.

## Motivation

**"I wish I had a stackoverflow page for the error this CLI just returned"**.

This is my attempt at making it easier for developers to improve the experience of their users during troubleshooting.

## 🚀 Get Started

<p align="center">
<img src="./docs/demo.gif">
</p>

Create an empty `errors.yaml` file in the same directory as the main.go.

```shell
touch error.yaml
```
Add in-code annotations for the errors the application's users might see.
```go
    //go:generate go run main.go generate -o errors.yaml --log-level none
    //go:generate go run main.go generate -o ./docs --log-level none
    package main
    
    import (
        "errors"
        "log"
    
        fyi "github.com/tfadeyi/errors"
    )

    //go:embed errors.yaml
    var errorsYAML []byte

    // @fyi name example
    // @fyi base_url https://example.github.io
    // @fyi version v0.0.1
    // @fyi description example error manifests
    func main() {
		fyi.SetSpecification(errorsYAML)
        err := doSomething()
        if err != nil {
            log.Fatal(err)
        }
    }
    func doSomething() error {
        // @fyi.error code clean_artefacts_error
        // @fyi.error title Error Removing Previous Artefacts
        // @fyi.error short The tool has failed to delete the artefacts from the previous execution.
        err := errors.New("something")
        return fyi.Error(err, "error_something_code")
    }
```

```shell
go generate ./... # will generate the errors documentation and manifests 
```

When returning an error the application will now return an error with the additional context described in the annotations:

```text
$ go run main.go
Error: [something]
The tool has failed to delete the artefacts from the previous execution.
for additional info check https://github.com/tfadeyi/my-app/my-app/error_something_code
```

## 📚 Features

- Easy to embed in existing Go applications with **go:generate** and **embed**.
- Generate markdown error documentation, with both **default** and **custom** markdown templates.
- Simple integration with static page generators.
- Simple **library** interface for an easy adoption.
- Error **embedded auto-generated URL** linking to error documentation.

## 🔽 Installation

### Library


```sh
 go get -u github.com/tfadeyi/errors@latest
```

### CLI

<details>
<summary>Linux</summary>

> The recommended way to install error.fyi CLI (errctl) is via the installation script:

> ```sh
> curl -sfL https://raw.githubusercontent.com/slosive/sloscribe/main/install.sh | sh -
> ```

</details>

<details>
<summary>MacOS</summary>

> The recommended way to install error.fyi CLI (errctl) is via the installation script:

> ```sh
> curl -sfL https://raw.githubusercontent.com/slosive/sloscribe/main/install.sh | sh -
> ```

</details>

## Try it!

### Nix

```shell
nix develop github:tfadeyi/errors
```

## 🖥️  Usage

```text
Usage:
  errctl generate [flags]

Flags:
      --error-template string   
  -f, --file string             Source code file to parse
      --format string           Output format (yaml,markdown) (default "yaml")
  -h, --help                    help for generate
  -d, --include strings         Comma separated list of directories to be parses by the tool (default [/home/oluwole/go/src/github.com/tfadeyi/errors/cmd])
      --info-template string    
  -l, --language string         Target source code language (default "go")
  -o, --output string           Target output file or directory to store the generated output

Global Flags:
      --log-level string   Only log messages with the given severity or above. One of: [none, debug, info, warn], errors will always be printed (default "info")
```

## License
error.fyi is released under the [MIT](./LICENSE) License.
