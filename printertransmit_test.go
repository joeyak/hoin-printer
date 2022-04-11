package hoin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransmitPrinterStatus(t *testing.T) {
	_, printer := newPrinter()

	_, err := printer.TransmitPrinterStatus()

	assert.NoError(t, err)
}

func TestTransmitOfflineStatus(t *testing.T) {
	_, printer := newPrinter()

	_, err := printer.TransmitOfflineStatus()

	assert.NoError(t, err)
}

func TestTransmitErrorStatus(t *testing.T) {
	_, printer := newPrinter()

	_, err := printer.TransmitErrorStatus()

	assert.NoError(t, err)
}

func TestTransmitPaperSensorStatus(t *testing.T) {
	_, printer := newPrinter()

	_, err := printer.TransmitPaperSensorStatus()

	assert.NoError(t, err)
}
