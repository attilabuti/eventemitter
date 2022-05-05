# Event Emitter

[![build status](https://img.shields.io/github/workflow/status/attilabuti/eventemitter/CI/main?style=flat-square)](https://github.com/attilabuti/eventemitter/actions)
[![codecov](https://img.shields.io/codecov/c/github/attilabuti/eventemitter?style=flat-square)](https://codecov.io/gh/attilabuti/eventemitter)
[![Go Report Card](https://goreportcard.com/badge/github.com/attilabuti/eventemitter?style=flat-square)](https://goreportcard.com/report/github.com/attilabuti/eventemitter)
[![Go Reference](https://pkg.go.dev/badge/github.com/attilabuti/eventemitter.svg)](https://pkg.go.dev/github.com/attilabuti/eventemitter)
[![license](https://img.shields.io/github/license/attilabuti/eventemitter?style=flat-square)](https://raw.githubusercontent.com/attilabuti/eventemitter/main/LICENSE)

Simple Event Emitter for Go Programming Language 1.18+.

## Installation

```bash
$ go get github.com/attilabuti/eventemitter@latest
```

## Usage

**For more information, please see the [Package Docs](https://pkg.go.dev/github.com/attilabuti/eventemitter).**

```go
package main

import (
	"fmt"

	"github.com/attilabuti/eventemitter"
)

func main() {
    // Creating an instance.
    emitter := eventemitter.New()

    // Register an event listener.
    emitter.AddListener("test_event", func(name string) {
        fmt.Printf("Hello, %s!", name)
    })

    // Emit event sync.
    emitter.EmitSync("test_event", "World")

    // Remove event listener.
    emitter.RemoveListener("test_event")
}
```

## Issues

Submit the [issues](https://github.com/attilabuti/eventemitter/issues) if you find any bug or have any suggestion.

## Contribution

Fork the [repo](https://github.com/attilabuti/eventemitter) and submit pull requests.

## License

This project is licensed under the [MIT License](https://github.com/attilabuti/eventemitter/blob/main/LICENSE).