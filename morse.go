package hoin

import (
	"fmt"
	"strings"
	"time"
)

func stringToMorse(s string) []int {
	var data []int
	for _, word := range strings.Split(strings.ToLower(s), " ") {
		for _, r := range word {
			var timings []int
			switch r {
			case 'a':
				timings = []int{1, 3}
			case 'b':
				timings = []int{3, 1, 1, 1}
			case 'c':
				timings = []int{3, 1, 3, 1}
			case 'd':
				timings = []int{3, 1, 1}
			case 'e':
				timings = []int{1}
			case 'f':
				timings = []int{1, 1, 3, 1}
			case 'g':
				timings = []int{3, 3, 1}
			case 'h':
				timings = []int{1, 1, 1, 1}
			case 'i':
				timings = []int{1, 1}
			case 'j':
				timings = []int{1, 3, 3, 3}
			case 'k':
				timings = []int{3, 1, 3}
			case 'l':
				timings = []int{1, 3, 1, 1}
			case 'm':
				timings = []int{3, 3}
			case 'n':
				timings = []int{3, 1}
			case 'o':
				timings = []int{3, 3, 3}
			case 'p':
				timings = []int{1, 3, 3, 1}
			case 'q':
				timings = []int{3, 3, 1, 3}
			case 'r':
				timings = []int{1, 3, 1}
			case 's':
				timings = []int{1, 1, 1}
			case 't':
				timings = []int{3}
			case 'u':
				timings = []int{1, 1, 3}
			case 'v':
				timings = []int{1, 1, 1, 3}
			case 'w':
				timings = []int{1, 3, 3}
			case 'x':
				timings = []int{3, 1, 1, 3}
			case 'y':
				timings = []int{3, 1, 3, 3}
			case 'z':
				timings = []int{3, 3, 1, 1}
			case '0':
				timings = []int{3, 3, 3, 3, 3}
			case '1':
				timings = []int{1, 3, 3, 3, 3}
			case '2':
				timings = []int{1, 1, 3, 3, 3}
			case '3':
				timings = []int{1, 1, 1, 3, 3}
			case '4':
				timings = []int{1, 1, 1, 1, 3}
			case '5':
				timings = []int{1, 1, 1, 1, 1}
			case '6':
				timings = []int{3, 1, 1, 1, 1}
			case '7':
				timings = []int{3, 3, 1, 1, 1}
			case '8':
				timings = []int{3, 3, 3, 1, 1}
			case '9':
				timings = []int{3, 3, 3, 3, 1}
			default:
				continue
			}
			for _, t := range timings {
				if t == 0 {
					data = append(data, 0, 0)
					continue
				}
				data = append(data, t, 0)
			}
			data = append(data, make([]int, 3)...)
		}
		data = append(data, make([]int, 7)...)
	}
	return data
}

// Morse beeps out the message in morse code
func (p Printer) Morse(message string) error {
	data := stringToMorse(message)
	for _, t := range data {
		if t != 0 {
			err := p.Beep(1, t)
			if err != nil {
				return fmt.Errorf("could not send morse code beeps")
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
