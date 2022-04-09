package hoin

import (
	"fmt"
	"io"
)

const (
	HT  = 0x09
	LF  = 0x0A
	CR  = 0x0D
	ESC = 0x1B
	DLE = 0x10
	EOT = 0x04
)

type Printer struct {
	dst io.Writer
}

func NewPrinter(dst io.ReadWriter) Printer {
	return Printer{
		dst: dst,
	}
}

func (p *Printer) WriteRaw(b []byte) error {
	_, err := p.dst.Write(b)
	if err != nil {
		return fmt.Errorf("could not write to printer: %w", err)
	}
	return nil
}

func (p *Printer) Initialize() error {
	err := p.WriteRaw([]byte{ESC, '@'})
	if err != nil {
		return fmt.Errorf("could not initialize printer: %w", err)
	}
	return nil
}

func (p *Printer) Print(a ...any) error {
	err := p.WriteRaw([]byte(fmt.Sprint(a...)))
	if err != nil {
		return fmt.Errorf("could not print %v: %w", a, err)
	}
	return nil
}

func (p *Printer) Println(a ...any) error {
	return p.Print(fmt.Sprint(a...) + string(rune(LF)))
}

func (p *Printer) Printf(format string, a ...any) error {
	return p.Print(fmt.Sprintf(format, a...))
}

func (p *Printer) HT() error {
	return p.WriteRaw([]byte{HT})
}

func (p *Printer) LF() error {
	return p.WriteRaw([]byte{LF})
}

func (p *Printer) CR() error {
	return p.WriteRaw([]byte{CR})
}
