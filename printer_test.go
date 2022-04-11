package hoin_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/joeyak/hoin-printer"
	"github.com/stretchr/testify/assert"
)

func convertIntsToBytes(a []int) []byte {
	var data []byte
	for _, b := range a {
		data = append(data, byte(b))
	}
	return data
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

func FuzzWrite(f *testing.F) {
	f.Add([]byte("Test"))
	f.Fuzz(func(t *testing.T, b []byte) {
		buffer, printer := newPrinter()

		n, err := printer.Write(b)

		assert.NoError(t, err)
		assert.Equal(t, len(b), n)
		assert.Equal(t, string(b), buffer.String())
	})
}

func FuzzRead(f *testing.F) {
	f.Add([]byte("Test"), 0)
	f.Fuzz(func(t *testing.T, b []byte, n int) {
		buffer, printer := newPrinter()
		buffer.Write(b)

		data := make([]byte, len(b)+n)
		n, err := printer.Read(data)

		assert.NoError(t, err)
		assert.Equal(t, len(b), n)
		assert.Equal(t, b, data[:n])
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
	buffer, printer := newPrinter()

	err := printer.CutFeed(5)

	assert.NoError(t, err)
	assert.Equal(t, "\x1DV\x42\x05", buffer.String())
}

func TestResetLineSpacing(t *testing.T) {
	buffer, printer := newPrinter()

	err := printer.ResetLineSpacing()

	assert.NoError(t, err)
	assert.Equal(t, "\x1B2", buffer.String())
}

func TestSetLineSpacing(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
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
	testCases := []struct {
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

	for _, tc := range testCases {
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

func TestFeedLines(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d:%t", tc.units, tc.err), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.FeedLines(tc.units)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte{0x1B, 'd', byte(tc.units)}, buffer.Bytes()[buffer.Len()-3:])
			}
		})
	}
}

func TestSetHTPositions(t *testing.T) {
	testCases := [][]int{
		{},
		{1},
		{4},
		{4, 4, 4},
	}

	for _, b := range testCases {
		t.Run(fmt.Sprint(b), func(t *testing.T) {
			expected := convertIntsToBytes(append(append([]int{0x1B, 'D'}, b...), 0))
			buffer, printer := newPrinter()

			err := printer.SetHT(b...)

			assert.NoError(t, err)
			assert.EqualValues(t, expected, buffer.Bytes())
		})
	}
}

func TestSetHTPositionsError(t *testing.T) {
	testCases := [][]int{
		{-1},
		{256},
		make([]int, 256),
	}

	for _, b := range testCases {
		t.Run(fmt.Sprint(b), func(t *testing.T) {
			_, printer := newPrinter()

			err := printer.SetHT(b...)

			assert.Error(t, err)
		})
	}
}

func TestSetBold(t *testing.T) {
	testCases := []struct {
		input  bool
		output byte
	}{
		{false, 0},
		{true, 1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.input), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.SetBold(tc.input)

			assert.NoError(t, err)
			assert.Equal(t, []byte{0x1B, 'E', tc.output}, buffer.Bytes())
		})
	}
}

func TestSetFont(t *testing.T) {
	for _, i := range []int{0, 1} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.SetFont(i)

			assert.NoError(t, err)
			assert.Equal(t, []byte{0x1B, 'M', byte(i)}, buffer.Bytes())
		})
	}
}

func TestSetFontErr(t *testing.T) {
	for _, i := range []int{-1, 2, 3} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			_, printer := newPrinter()

			err := printer.SetFont(i)

			assert.Error(t, err)
		})
	}
}

func TestSetRotate90(t *testing.T) {
	testCases := []struct {
		input  bool
		output byte
	}{
		{false, 0},
		{true, 1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.input), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.SetRotate90(tc.input)

			assert.NoError(t, err)
			assert.Equal(t, []byte{0x1B, 'V', tc.output}, buffer.Bytes())
		})
	}
}

func TestBeep(t *testing.T) {
	testCases := []struct {
		n, t int
		err  bool
	}{
		{0, 0, true},
		{1, 0, true},
		{0, 1, true},
		{1, 1, false},
		{9, 9, false},
		{10, 10, true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%d-%t", tc.n, tc.t, tc.err), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.Beep(tc.n, tc.t)

			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, []byte{0x1B, 'B', byte(tc.n), byte(tc.t)}, buffer.Bytes()[buffer.Len()-4:])
			}
		})
	}
}

func TestJustify(t *testing.T) {
	for _, j := range []hoin.Justification{hoin.Left, hoin.Center, hoin.Right} {
		t.Run(fmt.Sprint(j), func(t *testing.T) {
			buffer, printer := newPrinter()

			err := printer.Justify(j)

			assert.NoError(t, err)
			assert.Equal(t, []byte{0x1B, 'a', byte(j)}, buffer.Bytes())
		})
	}
}
