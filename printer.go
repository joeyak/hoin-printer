package hoin

import (
	"fmt"
	"io"
)

const (
	HT  = 0x09
	LF  = 0x0A
	CR  = 0x0D
	GS  = 0x1D
	ESC = 0x1B
	DLE = 0x10
	EOT = 0x04
)

func checkRange(n, min, max int, info string) error {
	if n < min || max < n {
		return fmt.Errorf("%s must be between %d and %d", info, min, max)
	}
	return nil
}

type Printer struct {
	dst io.ReadWriter
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
		return fmt.Errorf("could not print %q: %w", a, err)
	}
	return nil
}

func (p *Printer) Println(a ...any) error {
	return p.Print(fmt.Sprint(a...) + string(rune(LF)))
}

func (p *Printer) Printf(format string, a ...any) error {
	return p.Print(fmt.Sprintf(format, a...))
}

// HT moves the print position to the next horizontal tab position
func (p *Printer) HT() error {
	return p.WriteRaw([]byte{HT})
}

// LF prints the data in the print buffer and feeds one line
func (p *Printer) LF() error {
	return p.WriteRaw([]byte{LF})
}

// CR prints and does a carriage return
func (p *Printer) CR() error {
	return p.WriteRaw([]byte{CR})
}

// Cut cuts the paper
func (p *Printer) Cut() error {
	return p.WriteRaw([]byte{GS, 'V', 0})
}

// ResetLineSpacing sets the spacing to the default which
// is 1/6-inch lines (approx. 4.23mm)
func (p *Printer) ResetLineSpacing() error {
	return p.WriteRaw([]byte{ESC, '2'})
}

// SetLineSpacing sets the line spacing to n * v/h motion units in inches
func (p *Printer) SetLineSpacing(n int) error {
	if err := checkRange(n, 0, 255, "n"); err != nil {
		return fmt.Errorf("could not set line spacing: %w", err)
	}
	return p.WriteRaw([]byte{ESC, '3', byte(n)})
}

// Feed feeds the paper n units
func (p *Printer) Feed(n int) error {
	if err := checkRange(n, 0, 255, "n"); err != nil {
		return fmt.Errorf("could not feed paper: %w", err)
	}
	return p.WriteRaw([]byte{ESC, 'J', byte(n)})
}

// FeedLines feeds the paper n lines
func (p *Printer) FeedLines(n int) error {
	if err := checkRange(n, 0, 255, "n"); err != nil {
		return fmt.Errorf("could not feed lines: %w", err)
	}
	return p.WriteRaw([]byte{ESC, 'd', byte(n)})
}
