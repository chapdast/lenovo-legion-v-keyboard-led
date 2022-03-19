package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	lvl "github.com/chapdast/legion-v-keyboard-led"
	"github.com/google/gousb"
)

var (
	effectType string
	waveDir    string
	brightness int
	speed      int
	colors     string
	vpid       string
)
var (
	cmdEffectType = `Keyboard LED light effect:
	static 
	breath
	wave
	hue
	off`
	cmdWaveDir = `wave effect direction
	rtl
	ltr`
	cmdLEDBrightness = `brightness of LEDs
	0-off
	1
	2-brightest`
	cmdEffectSpeed = `effect speed
	1-slowest
	2
	3
	4-fastest`
	cmdColors = `effect color
	eg.: 00ff00 | ff00ff,00fff`
	cmdVPID = `Vendor and Product ID in VID:PID format.`
)

func init() {
	flag.StringVar(&effectType, "effect", "off", cmdEffectType)
	flag.StringVar(&waveDir, "dir", "ltr", cmdWaveDir)
	flag.IntVar(&brightness, "bright", 2, cmdLEDBrightness)
	flag.IntVar(&speed, "speed", 1, cmdEffectSpeed)
	flag.StringVar(&colors, "colors", "", cmdColors)
	flag.StringVar(&vpid, "vpid", "048d:c955", cmdVPID)

	flag.StringVar(&effectType, "e", "off", cmdEffectType)
	flag.StringVar(&waveDir, "d", "ltr", cmdWaveDir)
	flag.IntVar(&brightness, "b", 2, cmdLEDBrightness)
	flag.IntVar(&speed, "s", 1, cmdEffectSpeed)
	flag.StringVar(&colors, "c", "", cmdColors)

}

func main() {

	flag.Parse()

	log.Printf("flags:\nET:%s,\nSP:%d,\nBR:%d,\nDR:%s,\nCS:%s\n", effectType, brightness, speed, waveDir, colors)

	lk := lvl.New()

	switch strings.ToLower(effectType) {
	case "static":
		list, err := lvl.ColorFromString(colors)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		if err := lk.Static(lvl.EffectSpeed(speed), lvl.Brightness(brightness), list[0]); err != nil {
			log.Fatalln(err)
		}
	case "breath":
		list, err := lvl.ColorFromString(colors)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		if err := lk.Breath(lvl.EffectSpeed(speed), lvl.Brightness(brightness), list...); err != nil {
			log.Fatalln(err)
		}
	case "hue":
		if err := lk.HUE(lvl.EffectSpeed(speed), lvl.Brightness(brightness)); err != nil {
			log.Fatalln(err)
		}
	case "wave":
		var dir lvl.Direction = lvl.Def
		if strings.ToLower(waveDir) == "ltr" {
			dir = lvl.LTR
		} else {
			dir = lvl.RTL
		}
		if err := lk.Wave(lvl.EffectSpeed(speed), lvl.Brightness(brightness), dir); err != nil {
			log.Fatalln(err)
		}
	case "off":
		if err := lk.Off(); err != nil {
			log.Fatalln(err)
		}
	default:
		fmt.Printf("unknown effect type: %q\n", effectType)
		os.Exit(1)
	}

	ctx := gousb.NewContext()
	defer ctx.Close()

	// dev, err := getDevice(ctx)
	vpID := strings.Split(vpid, ":")
	vIDtmp, err := strconv.ParseUint(vpID[0], 16, 16)
	if err != nil {
		log.Fatalln("error parsing vendor id")
	}
	pIDtmp, err := strconv.ParseUint(vpID[1], 16, 16)
	if err != nil {
		log.Fatalln("error parsing product id")
	}
	vID := uint16(vIDtmp) //0x048d
	pID := uint16(pIDtmp) //0xc955

	dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(vID), gousb.ID(pID))
	if err != nil {
		if errors.Is(err, gousb.ErrorAccess) {
			fmt.Println(`Add udev rule as "/etc/udev/rules.d/10-kblight.rules" if you want control light as user
			SUBSYSTEM=="usb", ATTR{idVendor}=="048d", ATTR{idProduct}=="c965", MODE="0666"
			`)
			fmt.Println(ErrPermissionDenied)
			os.Exit(1)
		}
		fmt.Printf("error located device: %s\n", err)
		os.Exit(1)
	}
	if dev == nil {
		fmt.Println("can not get device make sure VendorID:ProducID is correct, current is % x:% x", vID, pID)
		os.Exit(1)
	}
	defer dev.Close()
	// Device need to get detached before sending any commands
	if err := dev.SetAutoDetach(true); err != nil {
		fmt.Printf("failed to detach device from kernel, %s\n", err)
		os.Exit(1)
	}

	data := lk.Data()
	log.Printf("data: %d, % x\n", len(data), data)
	fmt.Println(lk)
	c, err := dev.Control(0x21, 0x9, 0x03CC, 0x00, data)
	if err != nil {
		if errors.Is(err, gousb.ErrorBusy) {
			fmt.Println("error send command: device is busy")
			os.Exit(1)
		}
		fmt.Printf("error send command: %s\n", err)
		os.Exit(1)
	}
	if c != len(data) {
		fmt.Println(ErrInvalidDataLength)
		os.Exit(1)
	}
}

var (
	ErrInvalidDataLength = errors.New("invalid length of data written to device")
	ErrNoDeviceFound     = errors.New("no lenovo keyboard device found")
	ErrPermissionDenied  = errors.New("permission denied")
)

// func sendCommand(dev *gousb.Device, data []byte) error {

// }

// func getDevice(ctx *gousb.Context) (dev *gousb.Device, err error) {
// 	vID := 0x048d
// 	pIDs := []uint16{0xc955, 0x965}
// 	for _, pID := range pIDs {
// 		dev, err = ctx.OpenDeviceWithVIDPID(gousb.ID(vID), gousb.ID(pID))
// 		if err != nil {
// 			if errors.Is(err, gousb.ErrorAccess) {
// 				return nil, ErrPermissionDenied
// 			}
// 			log.Printf("error located device: %s\n", err)
// 		}
// 		if dev != nil {
// 			return dev, nil
// 		}
// 	}
// 	if dev == nil {
// 		return nil, ErrNoDeviceFound
// 	}
// 	return dev, err
// }
