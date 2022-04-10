package hoin_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/joeyak/hoin-printer"
	"github.com/stretchr/testify/assert"
)

var defaultUnitsTestCase = []struct {
	units int
	err   bool
}{
	{-1, true},
	{0, false},
	{1, false},
	{100, false},
	{255, false},
	{256, true},
}

func newPrinter() (*bytes.Buffer, hoin.Printer) {
	buffer := &bytes.Buffer{}
	return buffer, hoin.NewPrinter(buffer)
}

func TestHT(t *testing.T) {
	buffer, printer := newPrinter()

	err := printer.HT()

	assert.NoError(t, err)
	assert.Equal(t, "\x09", buffer.String())
}

func TestLF(t *testing.T) {
	buffer, printer := newPrinter()

	err := printer.LF()

	assert.NoError(t, err)
	assert.Equal(t, "\x0A", buffer.String())
}

func TestCR(t *testing.T) {
	buffer, printer := newPrinter()

	err := printer.CR()

	assert.NoError(t, err)
	assert.Equal(t, "\x0D", buffer.String())
}

func TestInitialize(t *testing.T) {
	buffer, printer := newPrinter()

	err := printer.Initialize()

	assert.NoError(t, err)
	assert.Equal(t, []byte{hoin.ESC, '@'}, buffer.Bytes())
}

func FuzzWriteRaw(f *testing.F) {
	f.Add([]byte("Test"))
	f.Fuzz(func(t *testing.T, b []byte) {
		buffer, printer := newPrinter()

		err := printer.WriteRaw(b)

		assert.NoError(t, err)
		assert.Equal(t, string(b), buffer.String())
	})
}

func FuzzPrintln(f *testing.F) {
	f.Add("Test")
	f.Fuzz(func(t *testing.T, s string) {
		buffer, printer := newPrinter()

		err := printer.Println(s)

		assert.NoError(t, err)
		assert.Equal(t, s+"\x0A", buffer.String())
	})
}

func FuzzPrintf(f *testing.F) {
	f.Add("Test")
	f.Fuzz(func(t *testing.T, s string) {
		a := []interface{}{1, 2}
		format := fmt.Sprintf("%s %%d %%d", s)
		expected := fmt.Sprintf(format, a...)

		buffer, printer := newPrinter()

		err := printer.Printf(format, a...)

		assert.NoError(t, err)
		assert.Equal(t, expected, buffer.String())
	})
}

func FuzzPrint(f *testing.F) {
	f.Add("Test")
	f.Fuzz(func(t *testing.T, s string) {
		buffer, printer := newPrinter()

		err := printer.Print(s)

		assert.NoError(t, err)
		assert.Equal(t, s, buffer.String())
	})
}

func TestCut(t *testing.T) {
	buffer, printer := newPrinter()

	err := printer.Cut()

	assert.NoError(t, err)
	assert.Equal(t, "\x1DV\x00", buffer.String())
}

func TestCutFeed(t *testing.T) {
	for _, tc := range defaultUnitsTestCase {
		t.Run(fmt.Sprintf("%d:%t", tc.units, tc.err), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.CutFeed(tc.units)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte{0x1D, 'V', 0x00, byte(tc.units)}, buffer.Bytes()[buffer.Len()-4:])
			}
		})
	}
}

func TestResetLineSpacing(t *testing.T) {
	buffer, printer := newPrinter()

	err := printer.ResetLineSpacing()

	assert.NoError(t, err)
	assert.Equal(t, "\x1B2", buffer.String())
}

func TestSetLineSpacing(t *testing.T) {
	for _, tc := range defaultUnitsTestCase {
		t.Run(fmt.Sprintf("%d:%t", tc.units, tc.err), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.SetLineSpacing(tc.units)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte{0x1B, '3', byte(tc.units)}, buffer.Bytes())
			}
		})
	}
}

func TestFeed(t *testing.T) {
	for _, tc := range defaultUnitsTestCase {
		t.Run(fmt.Sprintf("%d:%t", tc.units, tc.err), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.Feed(tc.units)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte{0x1B, 'J', byte(tc.units)}, buffer.Bytes()[buffer.Len()-3:])
			}
		})
	}
}
