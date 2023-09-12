package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/gousb"

	lvl "github.com/chapdast/legion-v-keyboard-led"
)

var (
	effectType string
	waveDir    string
	brightness int
	speed      int
	colors     string
	vpid       string
	debug      bool
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

var (
	permitHelpFormat = `Add udev rule as "/etc/udev/rules.d/10-kblight.rules" if you want control light as user
	SUBSYSTEM=="usb", ATTR{idVendor}=="%04x", ATTR{idProduct}=="%04x", MODE="0666"
	
	then run
	
	sudo udevadm control --reload-rules && sudo udevadm trigger
	`
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

	flag.BoolVar(&debug, "debug", false, "show debug messages")

}

func main() {

	flag.Parse()
	if debug {
		fmt.Printf("\n-COMMAND------------\n")
		fmt.Printf("flags:\n\tEffect Type: %s,\n\tSpeed: %d,\n\tBrightness: %d,\n\tDirection: %s,\n\tRGB: %s\n", 
		effectType, brightness, speed, waveDir, colors)
	}
	lk := lvl.New()

	switch strings.ToLower(effectType) {
	case "static":
		list, err := lvl.ColorFromString(colors)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}
		if err := lk.Static(lvl.EffectSpeed(speed), lvl.Brightness(brightness), list...); err != nil {
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

	if debug {
		ctx.Debug(1)
	}
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
			fmt.Printf(permitHelpFormat, vID, pID)
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
	// // Device need to get detached before sending any commands
	// if err := dev.SetAutoDetach(true); err != nil {
	// 	fmt.Printf("failed to detach device from kernel, %s\n", err)
	// 	os.Exit(1)
	// }

	if err := dev.Reset(); err !=nil {
		fmt.Printf("can not reset device, %s\n", err)
		os.Exit(1)
	}
	if err := dev.SetAutoDetach(false); err !=nil {
		fmt.Printf("can not disable auto detach, %s\n", err)
		os.Exit(1)
	}

	if debug {
		fmt.Printf("\n-DEVICE-------------\n")
		for i := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 10} {
			des, err := dev.GetStringDescriptor(i)
			if err != nil {
				continue
			}
			fmt.Printf("\tDescriptor  %d: %s\n", i, des)
		}
		fmt.Printf("\n-DEVICE-CONFIG------\n")
		num, err := dev.ActiveConfigNum()
		fmt.Printf("\tActive Config: %v, err: %v\n", num, err)
		fmt.Printf("\tDevice Info: %v\n",dev.String())
		

	}


	data := lk.Data()

	if debug {
		fmt.Printf("\n-PAYLOAD-------------\n")
		fmt.Printf("len:%d\n% x\n", len(data), data)
		fmt.Printf("\n---------------------\n")
		fmt.Println(lk)
		fmt.Printf("-----------------------\n")
	}

	c, err := dev.Control(0x21, 0x9, 0x03cc, 0x00, data)
	if err != nil {
		fmt.Printf("error: %q\n", err)
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
