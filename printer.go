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

func (p Printer) Write(b []byte) (int, error) {
	n, err := p.dst.Write(b)
	if err != nil {
		return n, fmt.Errorf("could not write to printer: %w", err)
	}
	return n, nil
}

func (p Printer) Read(b []byte) (int, error) {
	n, err := p.dst.Read(b)
	if err != nil {
		return n, fmt.Errorf("could not read from printer: %w", err)
	}
	return n, nil
}

func (p Printer) Initialize() error {
	_, err := p.Write([]byte{ESC, '@'})
	if err != nil {
		return fmt.Errorf("could not initialize printer: %w", err)
	}
	return nil
}

func (p Printer) Print(a ...any) error {
	_, err := p.Write([]byte(fmt.Sprint(a...)))
	if err != nil {
		return fmt.Errorf("could not print %q: %w", a, err)
	}
	return nil
}

func (p Printer) Println(a ...any) error {
	return p.Print(fmt.Sprint(a...) + string(rune(LF)))
}

func (p Printer) Printf(format string, a ...any) error {
	return p.Print(fmt.Sprintf(format, a...))
}

// HT moves the print position to the next horizontal tab position
//
// By default HT will do nothing if SetHT is not called with tab positions
func (p Printer) HT() error {
	_, err := p.Write([]byte{HT})
	if err != nil {
		return fmt.Errorf("could not send HT: %w", err)
	}
	return nil
}

// LF prints the data in the print buffer and feeds one line
func (p Printer) LF() error {
	_, err := p.Write([]byte{LF})
	if err != nil {
		return fmt.Errorf("could not send LF: %w", err)
	}
	return nil
}

// CR prints and does a carriage return
func (p Printer) CR() error {
	_, err := p.Write([]byte{CR})
	if err != nil {
		return fmt.Errorf("could not send CR: %w", err)
	}
	return nil
}

// Cut cuts the paper
func (p Printer) Cut() error {
	_, err := p.Write([]byte{GS, 'V', 0})
	if err != nil {
		return fmt.Errorf("could not cut paper: %w", err)
	}
	return nil
}

// CutFeed feeds the paper n units and then cuts it
func (p Printer) CutFeed(n int) error {
	errMsg := "could not feed and cut the paper: %w"

	err := checkRange(n, 0, 255, "n")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{GS, 'V', 66, byte(n)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}

// ResetLineSpacing sets the spacing to the default which
// is 1/6-inch lines (approx. 4.23mm)
func (p Printer) ResetLineSpacing() error {
	_, err := p.Write([]byte{ESC, '2'})
	if err != nil {
		return fmt.Errorf("could not reset line spacing: %w", err)
	}
	return nil
}

// SetLineSpacing sets the line spacing to n * v/h motion units in inches
func (p Printer) SetLineSpacing(n int) error {
	errMsg := "could not set line spacing: %w"

	err := checkRange(n, 0, 255, "n")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{ESC, '3', byte(n)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return err
}

// Feed feeds the paper n units
func (p Printer) Feed(n int) error {
	errMsg := "could not feed paper: %w"

	err := checkRange(n, 0, 255, "n")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{ESC, 'J', byte(n)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return err
}

// FeedLines feeds the paper n lines
func (p Printer) FeedLines(n int) error {
	errMsg := "could not feed lines: %w"

	err := checkRange(n, 0, 255, "n")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{ESC, 'd', byte(n)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return err
}

// SetHT sets the horizontal tab positions
//
// This command cancels previous SetHT commands
// Multiple positions can be set for tabbing
// A max of 32 positions can be set
// Calling SetHT with no argments resets the tab positions
func (p Printer) SetHT(positions ...int) error {
	errMsg := "could not set horizontal tab positions: %w"

	if len(positions) > 32 {
		return fmt.Errorf("more than 32 positions was set")
	}

	var data []byte
	for i, pos := range positions {
		err := checkRange(pos, 1, 255, fmt.Sprintf("position %d", i))
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		data = append(data, byte(pos))
	}

	data = append([]byte{ESC, 'D'}, data...)
	data = append(data, 0)

	_, err := p.Write(data)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

// SetBold turns emphasized mode on or off
func (p Printer) SetBold(b bool) error {
	var bb byte = 0
	if b {
		bb = 1
	}

	_, err := p.Write([]byte{ESC, 'E', bb})
	if err != nil {
		return fmt.Errorf("could not set bold: %w", err)
	}

	return nil
}

// SetFont changes the font
//
// n=0 selects font A
// n=1 selects font B
func (p Printer) SetFont(n int) error {
	errMsg := "could not set font: %w"

	if !(n == 0 || n == 1) {
		return fmt.Errorf(errMsg, fmt.Errorf("n must be 0 or 1"))
	}

	_, err := p.Write([]byte{ESC, 'M', byte(n)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
