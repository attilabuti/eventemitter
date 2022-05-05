# Event Emitter

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