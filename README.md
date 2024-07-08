# hoin-printer

[![Go Report Card](https://goreportcard.com/badge/github.com/joeyak/hoin-printer)](https://goreportcard.com/report/github.com/joeyak/hoin-printer)
![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)
[![GoDoc](https://godoc.org/github.com/joeyak/hoin-printer?status.svg)](https://godoc.org/github.com/joeyak/hoin-printer)
![tests](https://github.com/joeyak/hoin-printer/actions/workflows/main.yaml/badge.svg)

> [!IMPORTANT]
> This repo has been moved to [joeyak/go-escpos](https://github.com/joeyak/go-escpos) to match the fact that these commands work on other ESC/POS Thermal Printers

This is a package for writing to a HOIN POS-80-Series Thermal Printer

Connect to the printer with an io.ReadWriter and then send commands

```go
package main

import (
	"fmt"
	"net"

	"github.com/joeyak/hoin-printer"
)

func main() {

	conn, err := net.Dial("tcp", "192.168.1.23:9100")
	if err != nil {
		fmt.Println("unable to dial:", err)
		return
	}
	defer conn.Close()

	printer := hoin.NewPrinter(conn)

	for i := 0; i < 5; i++ {
		printer.Println("Hello World!")
	}

	printer.FeedLines(5)
	printer.Cut()
}
```

## Testing

What? Did I hear you ask for testing? You think we make useless mocks that only tests our assumptions about the hoin printer instead of REAL **HONEST** ***GOOD*** boots on the ground testing.

Run `go run ./cmd/test-printer/` to print out our test program.

Really, how are we supposed to tests without a firmware dump? Total incongruity.

Also the test program assumes some things will work line printing and the such, cause how can we test functions without that. It'd be obvious if nothing prints. The goal is to test all the extra functions like horizontal tabbing, justifications, images, etc.
