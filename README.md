# hoin-printer

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
