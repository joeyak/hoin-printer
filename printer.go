package hoin

import (
	"fmt"
	"image"
	"image/color"
	"io"
)

const (
	HT  = 0x09
	LF  = 0x0A
	CR  = 0x0D
	GS  = 0x1D
	ESC = 0x1B
	DLE = 0x10
)

type Justification byte

const (
	LeftJustify Justification = iota
	CenterJustify
	RightJustify
)

type HRIPosition byte

const (
	HRINone HRIPosition = iota
	HRIAbove
	HRIBelow
	HRIBoth
)

// Density represents the DPI to use when printing images.
type Density byte

const (
	// SingleDensity is 90dpi
	SingleDensity Density = iota
	// DoubleDensity is 180dpi
	DoubleDensity
)

type BarCode byte

const (
	BcUPCA BarCode = iota
	BcUPCE
	BcJAN13
	BcJAN8
	BcCODE39
	BcITF
	BcCODABAR
	BcCODE93
	BcCODE123
)

var (
	lengthBarcodes = []BarCode{BcCODE93, BcCODE123}
	allBarcodes    = append(lengthBarcodes, BcUPCA, BcUPCE, BcJAN13, BcJAN8, BcCODE39, BcITF, BcCODABAR)
)

func inSlice[T ~byte](v T, s ...T) bool {
	for _, a := range s {
		if v == a {
			return true
		}
	}
	return false
}

func checkEnum[T ~byte](e T, enums ...T) error {
	if inSlice(e, enums...) {
		return nil
	}
	return fmt.Errorf("%v was not a valid choice from %v", e, enums)
}

func checkRange(n, min, max int, info string) error {
	if n < min || max < n {
		return fmt.Errorf("%s must be between %d and %d", info, min, max)
	}
	return nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
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

// Beep makes a beep sound n times for t duration
//
// Duration is dependent on the model. For the HOP-E802
// each duration is around 100ms
func (p Printer) Beep(n, t int) error {
	errMsg := "could not beep the printer: %w"

	err := checkRange(n, 1, 9, "n")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	err = checkRange(t, 1, 9, "t")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{ESC, 'B', byte(n), byte(t)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
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
	_, err := p.Write([]byte{ESC, 'E', boolToByte(b)})
	if err != nil {
		return fmt.Errorf("could not set bold to %t: %w", b, err)
	}
	return nil
}

// SetRotate90 turns on 90 clockwise rotation mode for the text
//
// When text is double-width or double-height the text will be mirrored
func (p Printer) SetRotate90(b bool) error {
	_, err := p.Write([]byte{ESC, 'V', boolToByte(b)})
	if err != nil {
		return fmt.Errorf("could not set bold to %t: %w", b, err)
	}
	return nil
}

// SetReversePrinting sets the white/black printing mode
//
// If b is true then it will print black text on white background
// If b is false then it will print white text on black background
func (p Printer) SetReversePrinting(b bool) error {
	_, err := p.Write([]byte{GS, 'B', boolToByte(b)})
	if err != nil {
		return fmt.Errorf("could not set reverse printing mode: %w", err)
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

// Justify sets the alignment to n
func (p Printer) Justify(j Justification) error {
	errMsg := "could not justify: %w"

	err := checkEnum(j, CenterJustify, LeftJustify, RightJustify)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{ESC, 'a', byte(j)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}

// PrintImage8 prints an image in the 8-bit row format.  In this format each
// row is 8 dots tall.
//
// The density selects the horizontal DPI of the image.  SingleDensity is
// 90dpi while DoubleDensity is 180dpi.  Vertical DPI is always 60dpi for
// 8-bit image data.
//
// No black and white conversion is performed on the provided image.  The
// image should be converted before calling this function.
func (p Printer) PrintImage8(img image.Image, density Density) error {
	imgRect := img.Bounds()
	var err error
	errMsg := "could not print 8 dot image: %w"

	err = checkEnum(density, SingleDensity, DoubleDensity)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	// 8 dot density (meta row is 8 dots tall)
	for y := 0; y < imgRect.Max.Y; y += 8 {
		row := []byte{}
		for x := 0; x < imgRect.Max.X; x++ {
			col := byte(0)

			for i := 0; i < 8; i++ {
				col <<= 1
				// Pad the bottom row to be white
				if y+i > imgRect.Max.Y {
					continue
				}
				c := color.GrayModel.Convert(img.At(x, y+i)).(color.Gray)
				if c.Y == 0 {
					col |= 1
				}
			}

			row = append(row, col)
		}

		data := []byte{ESC, '*', byte(density), byte(len(row)), byte(len(row) >> 8)}

		if err = p.SetLineSpacing(0); err != nil {
			return fmt.Errorf(errMsg, err)
		}
		if _, err = p.Write(append(data, row...)); err != nil {
			return fmt.Errorf(errMsg, err)
		}
		if err = p.LF(); err != nil {
			return fmt.Errorf(errMsg, err)
		}
	}

	return nil
}

// PrintImage24 prints an image in the 24-bit row format.  In this format each
// row is 24 dots tall.
//
// This works the same as PrintImage8() with the only difference being the DPI
// of the printed image.  SingleDensity is 90dpi while DoubleDensity is
// 180dpi.  Vertical DPI is always 180dpi for 24-bit image data.
func (p Printer) PrintImage24(img image.Image, density Density) error {
	imgRect := img.Bounds()
	var err error
	errMsg := "could not print 24 dot image: %w"

	err = checkEnum(density, SingleDensity, DoubleDensity)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	imgBytes := [][]byte{}

	// First convert the image data to the 1-bit data format.
	// 24 dot density (meta row is 24 dots tall (3 bytes))
	for y := 0; y < imgRect.Max.Y; y += 24 {
		metaRow := []byte{}
		for x := 0; x < imgRect.Max.X; x++ {

			for z := 0; z < 3; z++ {
				col := byte(0)
				for i := 0; i < 8; i++ {
					col <<= 1
					// Pad the bottom row to be white
					if (y+z*8)+i > imgRect.Max.Y {
						continue
					}

					c := color.GrayModel.Convert(img.At(x, (y+z*8)+i)).(color.Gray)
					if c.Y == 0 {
						col |= 1
					}
				}
				metaRow = append(metaRow, col)
			}

		}
		imgBytes = append(imgBytes, metaRow)
	}

	// Next send the data to the printer.
	command := []byte{ESC, 0x2A, byte(density + 32), byte(imgRect.Max.X), byte(imgRect.Max.X >> 8)}
	for _, row := range imgBytes {
		err = p.SetLineSpacing(0)
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		_, err = p.Write(append(command, row...))
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		err = p.LF()
		if err != nil {
			return fmt.Errorf(errMsg, err)
		}

		// If data is sent to fast it won't make it to the printer and will
		// stop printing part of the way through an image.  This will also
		// lose any commands sent after the image.  Sleeping for 35ms seems to
		// be the best balance between not printing and reducing banding.
		// time.Sleep(time.Millisecond * 35)

		// Leaving above there just in case below is a footgun, but below should
		// wait till print buffer is done to write more lines
		p.TransmitErrorStatus()
	}

	return nil
}

// SetHRIPosition sets the printing position of the HRI characters
// in relation to the barcode
func (p Printer) SetHRIPosition(hp HRIPosition) error {
	errMsg := "could not set HRI position: %w"

	err := checkEnum(hp, HRINone, HRIAbove, HRIBelow, HRIBoth)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{GS, 'H', byte(hp)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}
	return nil
}

// ResetBarCodeHeight sets the barcode height to 162
func (p Printer) ResetBarCodeHeight() error {
	err := p.SetBarCodeHeight(162)
	if err != nil {
		return fmt.Errorf("could not reset barcode height: %w", err)
	}
	return nil
}

// SetBarCodeHeight sets the barcode height in n dots
func (p Printer) SetBarCodeHeight(n int) error {
	errMsg := "could not set bar code height: %w"

	err := checkRange(n, 1, 255, "height")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	_, err = p.Write([]byte{GS, 'h', byte(n)})
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

func (p Printer) PrintBarCode(barcodeType BarCode, data ...byte) error {
	errMsg := "could not print bar code: %w"

	err := checkEnum(barcodeType, allBarcodes...)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	dataLength := len(data)

	// Check length
	var min, max int
	if barcodeType == BcUPCA {
		min, max = 11, 12
	}

	err = checkRange(dataLength, min, max, "data length")
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	msg := []byte{0x1D, 'k', byte(barcodeType)}

	// Add data
	if inSlice(barcodeType, lengthBarcodes...) {
		// length defined barcode
		msg = append(msg, byte(dataLength))
		msg = append(msg, data...)
	} else {
		// Null ending barcode
		msg = append(msg, data...)
		msg = append(msg, 0)
	}

	_, err = p.Write(msg)
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}
