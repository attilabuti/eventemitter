# Event Emitter

[![build status](https://img.shields.io/github/workflow/status/attilabuti/eventemitter/CI/main?style=flat-square)](https://github.com/attilabuti/eventemitter/actions)
[![codecov](https://img.shields.io/codecov/c/github/attilabuti/eventemitter?style=flat-square)](https://codecov.io/gh/attilabuti/eventemitter)
[![Go Report Card](https://goreportcard.com/badge/github.com/attilabuti/eventemitter/v2?style=flat-square)](https://goreportcard.com/report/github.com/attilabuti/eventemitter/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/attilabuti/eventemitter.svg)](https://pkg.go.dev/github.com/attilabuti/eventemitter/v2)
[![license](https://img.shields.io/github/license/attilabuti/eventemitter?style=flat-square)](https://raw.githubusercontent.com/attilabuti/eventemitter/main/LICENSE)

Simple Event Emitter for Go Programming Language 1.18+.

## Installation

```bash
$ go get github.com/attilabuti/eventemitter/v2@latest
```

## Usage

**For more information, please see the [Package Docs](https://pkg.go.dev/github.com/attilabuti/eventemitter/v2).**

```go
package main

import (
	"fmt"

	"github.com/attilabuti/eventemitter/v2"
)

func main() {
    // Creating an instance.
    emitter := eventemitter.New()

    // Event handler.
    event := func(name string) {
        fmt.Printf("Hello, %s!", name)
    }

    // Register an event listener.
    emitter.AddListener("test_event", event)

    // Emit event sync.
    emitter.EmitSync("test_event", "World")

    // Remove event listener.
    emitter.RemoveListener("test_event", event)
}
```

## Examples

### AddListener

```go
func main() {
	emitter := eventemitter.New()

    event := func(name string) {
        fmt.Printf("Hello, %s!", name)
    }

    emitter.AddListener("event", event)
    emitter.AddListener("event", &event)
    emitter.AddListener("event", func(name string) {
        fmt.Printf("Hello, %s!", name)
    })
}
```

### RemoveListener

```go
func main() {
	emitter := eventemitter.New()

    event := func(name string) {
        fmt.Printf("Hello, %s!", name)
    }

    emitter.AddListener("event", event)
    emitter.AddListener("event", &event)

    emitter.RemoveListener("event", event)
    emitter.RemoveListener("event", &event)
}
```

### Emit

```go
func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	emitter := eventemitter.New()

    emitter.AddListener("event", func(name string) {
        defer wg.Done()
        fmt.Printf("Hello, %s!", name)
    })

    emitter.Emit("event", "World")

    wg.Wait()
}
```

### EmitSync

```go
func main() {
	emitter := eventemitter.New()

    emitter.AddListener("event", func(name string) {
        fmt.Printf("Hello, %s!", name)
    })

    emitter.EmitSync("event", "World")
}
```

### RemoveAllListeners

```go
func main() {
	emitter := eventemitter.New()

    // Removes all listeners of the specified event.
    emitter.AddListener("event", func(){})
    emitter.AddListener("event", func(){})
    emitter.RemoveAllListeners("event")

    // Removes all listeners.
    emitter.AddListener("event1", func(){})
    emitter.AddListener("event2", func(){})
    emitter.RemoveAllListeners()
}
```

## Issues

Submit the [issues](https://github.com/attilabuti/eventemitter/issues) if you find any bug or have any suggestion.

## Contribution

Fork the [repo](https://github.com/attilabuti/eventemitter) and submit pull requests.

## License

This project is licensed under the [MIT License](https://github.com/attilabuti/eventemitter/blob/main/LICENSE).
