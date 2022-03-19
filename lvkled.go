package legion_v_keyboard_led

import (
	"errors"
)

var ErrNoColor = errors.New("no color specified")
var (
	ColorOff   = &Color{}
	ColorBlack = &Color{}
	ColorWhite = &Color{Red: 0xff, Green: 0xff, Blue: 0xff}
	ColorRed   = &Color{Red: 0xff}
	ColorGreen = &Color{Green: 0xff}
	ColorBlue  = &Color{Blue: 0xff}
)

type LKeyboard interface {
	Static(s EffectSpeed, b Brightness, c *Color) error
	Breath(s EffectSpeed, b Brightness, c ...*Color) error
	Wave(s EffectSpeed, b Brightness, rtl Direction) error
	HUE(s EffectSpeed, b Brightness) error
	Data() []byte

	Off() error
	// CUSTOM Control
	Manual(effect EffectType, speed EffectSpeed, bright Brightness, dir Direction, c ...*Color) error
}

var _ LKeyboard = &lKeyboard{}

type lKeyboard struct {
	dataPackage []byte
	effect      EffectType
	speed       EffectSpeed
	brightness  Brightness
	colors      []*Color
	waveRTL     Direction
}

type Direction int

const (
	Def = iota
	RTL
	LTR
)

type Color struct {
	Red   byte
	Green byte
	Blue  byte
}

func (c Color) RGB() []byte {
	return []byte{c.Red, c.Green, c.Blue}
}

func New() LKeyboard {
	return &lKeyboard{}
}

type Brightness byte

var (
	BrightnessLow  = Brightness(0x01)
	BrightnessHigh = Brightness(0x02)
)

type EffectType byte

var (
	EffectStatic = EffectType(0x01)
	EffectBreath = EffectType(0x03)
	EffectWave   = EffectType(0x04)
	EffectHue    = EffectType(0x06)
)

type EffectSpeed byte

var (
	EffectSpeedDefault = EffectSpeed(0x00)
	EffectSpeedSlowest = EffectSpeed(0x01)
	EffectSpeedSlow    = EffectSpeed(0x02)
	EffectSpeedFast    = EffectSpeed(0x03)
	EffectSpeedFastest = EffectSpeed(0x04)
)

type ControlBytes byte

var (
	ON  = ControlBytes(0x01)
	OFF = ControlBytes(0x00)
)

func (lk *lKeyboard) Static(s EffectSpeed, b Brightness, c *Color) error {
	return lk.Manual(EffectStatic, s, b, Def, c)
}
func (lk *lKeyboard) Breath(s EffectSpeed, b Brightness, c ...*Color) error {
	return lk.Manual(EffectBreath, s, b, Def, c...)
}
func (lk *lKeyboard) Wave(s EffectSpeed, b Brightness, dir Direction) error {
	return lk.Manual(EffectWave, s, b, dir, nil)
}

func (lk *lKeyboard) HUE(s EffectSpeed, b Brightness) error {
	return lk.Manual(EffectHue, s, b, Def, nil)
}

func (lk *lKeyboard) Data() []byte {
	return lk.dataPackage
}

func (lk *lKeyboard) Off() error {
	lk.effect = EffectStatic
	lk.brightness = 0
	return lk.Manual(EffectStatic, EffectSpeedFastest, 0, Def, ColorBlack)
}

func (lk *lKeyboard) Manual(effect EffectType, speed EffectSpeed, bright Brightness, dir Direction, c ...*Color) error {

	if (effect == EffectBreath || effect == EffectStatic) && len(c) == 0 {
		return ErrNoColor
	}
	lk.effect = effect
	lk.brightness = bright
	lk.speed = speed
	lk.colors = c
	lk.waveRTL = dir
	return lk.build()
}
