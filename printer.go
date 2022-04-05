package hoin

import (
	"fmt"
	"io"
)

const (
	HT  = "\x09"
	LF  = "\x0A"
	CR  = "\x0D"
	DLE = "\x10"
	EOT = "\x04"
)

type Printer struct {
	Writer io.Writer
	Width  int
}

func (p *Printer) write(b []byte) error {
	if p.Writer == nil {
		return fmt.Errorf("no writer was set")
	}

	_, err := p.Writer.Write(b)
	if err != nil {
		return fmt.Errorf("could not write to printer: %w", err)
	}

	return nil
}

func (p *Printer) Print(s string) error {
	err := p.write([]byte(s))
	if err != nil {
		return fmt.Errorf("could not print %q: %w", s, err)
	}
	return nil
}

func (p *Printer) Println(s string) error {
	return p.Print(s + LF)
}

func (p *Printer) Printf(format string, a ...any) error {
	return p.Print(fmt.Sprintf(format, a...))
}

func (p *Printer) HT() error {
	return p.Print(HT)
}

func (p *Printer) LF() error {
	return p.Print(LF)
}

func (p *Printer) CR() error {
	return p.Print(CR)
}
