package legion_v_keyboard_led

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ColorOff   = &Color{}
	ColorBlack = &Color{}
	ColorWhite = &Color{Red: 0xff, Green: 0xff, Blue: 0xff}
	ColorRed   = &Color{Red: 0xff}
	ColorGreen = &Color{Green: 0xff}
	ColorBlue  = &Color{Blue: 0xff}
)

type Color struct {
	Red   byte
	Green byte
	Blue  byte
}

func (c Color) String() string {
	return fmt.Sprintf("{r:%x, g:%x, b:%x}", c.Red, c.Green, c.Blue)
}
func (c Color) RGB() []byte {
	return []byte{c.Red, c.Green, c.Blue}
}

var ErrInvalidColor = errors.New("invalid color code")

func ColorFromString(cs string) ([]*Color, error) {
	clist := strings.Split(cs, ",")
	var list []*Color
	for _, c := range clist {
		if len(c) != 6 {
			return nil, ErrInvalidColor
		}
		red, err := strconv.ParseUint((c[0:2]), 16, 16)
		if err != nil {
			return nil, ErrInvalidColor
		}

		green, err := strconv.ParseUint((c[2:4]), 16, 16)
		if err != nil {
			return nil, ErrInvalidColor
		}

		blue, err := strconv.ParseUint((c[4:6]), 16, 16)
		if err != nil {
			return nil, ErrInvalidColor
		}

		list = append(list, &Color{Red: byte(red), Blue: byte(blue), Green: byte(green)})
	}
	return list, nil
}
