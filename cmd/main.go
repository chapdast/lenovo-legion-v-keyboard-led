package main

import (
	"errors"
	"log"

	lvl "github.com/chapdast/legion-v-keyboard-led"
	"github.com/google/gousb"
)

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	// dev, err := getDevice(ctx)

	vID := 0x048d
	pID := 0xc955

	dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(vID), gousb.ID(pID))
	if err != nil {
		if errors.Is(err, gousb.ErrorAccess) {
			log.Fatal(ErrPermissionDenied)
		}
		log.Fatalf("error located device: %s\n", err)
	}
	defer dev.Close()
	// Device need to get detached before sending any commands
	if err := dev.SetAutoDetach(true); err != nil {
		log.Fatalf("failed to detach device from kernel, %s\n", err)
	}

	lk := lvl.New()
	if err := lk.Breath(lvl.EffectSpeedDefault, lvl.BrightnessHigh, lvl.ColorGreen, lvl.ColorBlue); err != nil {
		log.Fatalf("error set effect, %s\n", err)
	}

	data := lk.Data()
	log.Printf("data: %d, % x", len(data), data)
	c, err := dev.Control(0x21, 0x9, 0x03CC, 0x00, data)
	if err != nil {
		if errors.Is(err, gousb.ErrorBusy) {
			log.Fatalln("error send command: device is busy")
		}
		log.Fatalf("error send command: %s\n", err)
	}
	if c != len(data) {
		log.Fatalln(ErrInvalidDataLength)
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
