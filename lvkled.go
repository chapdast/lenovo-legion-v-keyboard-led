package legion_v_keyboard_led

import (
	"errors"
	"fmt"
)

var ErrNoColor = errors.New("no color specified")

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
	dataPackage   []byte
	effect        EffectType
	speed         EffectSpeed
	brightness    Brightness
	colors        []*Color
	waveDirection Direction
}

type Direction int

func (d Direction) String() string {
	switch d {
	case LTR:
		return "LTR"
	case RTL:
		return "RTL"
	case Def:
		fallthrough
	default:
		return "NO-DIR"
	}
}

const (
	Def = iota
	RTL
	LTR
)

func New() LKeyboard {
	return &lKeyboard{}
}

type Brightness byte

func (b Brightness) String() string {
	switch b {
	case BrightnessLow:
		return "low"
	case BrightnessHigh:
		return "high"
	case BrightnessDefault:
		fallthrough
	default:
		return "off"
	}
}

var (
	BrightnessDefault = Brightness(0x00)
	BrightnessLow     = Brightness(0x01)
	BrightnessHigh    = Brightness(0x02)
)

type EffectType byte

func (et EffectType) String() string {
	switch et {
	case EffectStatic:
		return "static"
	case EffectBreath:
		return "breath"
	case EffectWave:
		return "wave"
	case EffectHue:
		return "hue"
	default:
		return "off"
	}
}

var (
	EffectStatic = EffectType(0x01)
	EffectBreath = EffectType(0x03)
	EffectWave   = EffectType(0x04)
	EffectHue    = EffectType(0x06)
)

type EffectSpeed byte

func (es EffectSpeed) String() string {
	switch es {
	case EffectSpeedSlowest:
		return "slowest"
	case EffectSpeedSlow:
		return "slow"
	case EffectSpeedFast:
		return "fast"
	case EffectSpeedFastest:
		return "fastest"
	case EffectSpeedDefault:
		fallthrough
	default:
		return "stoped"

	}
}

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
	lk.waveDirection = dir
	return lk.build()
}

func (lk lKeyboard) String() string {

	return fmt.Sprintf(`Configed with:
	Effect: %s,
	Speed: %s,
	Brightness: %s,
	Direction: %s,
	Colors: %s`, lk.effect, lk.speed, lk.brightness, lk.waveDirection, lk.colors)
}
