package legion_v_keyboard_led

func (lk *lKeyboard) header() {
	lk.dataPackage[0] = 0xCC
	lk.dataPackage[1] = 0x16
}

func (lk *lKeyboard) setEffect() {
	lk.dataPackage[2] = byte(lk.effect)
	lk.dataPackage[3] = byte(lk.speed)
	lk.dataPackage[4] = byte(lk.brightness)
}
func (lk *lKeyboard) setColor() error {
	if (lk.effect == EffectStatic || lk.effect == EffectBreath) && len(lk.colors) == 0 {
		return ErrNoColor
	}
	if lk.effect == EffectStatic {
		for i := 0; i < 4; i++ {
			lk.dataPackage[5+(i*3)+0] = lk.colors[0].Red
			lk.dataPackage[5+(i*3)+1] = lk.colors[0].Green
			lk.dataPackage[5+(i*3)+2] = lk.colors[0].Blue
		}
	} else if lk.effect == EffectBreath {
		for i := 0; i < 4; i++ {
			if i < len(lk.colors) {
				lk.dataPackage[5+(i*3)+0] = lk.colors[i].Red
				lk.dataPackage[5+(i*3)+1] = lk.colors[i].Green
				lk.dataPackage[5+(i*3)+2] = lk.colors[i].Blue
			}
		}
		//Fill Empty Colors
		if lk.effect == EffectBreath && len(lk.colors) != 4 {
			j := len(lk.colors)
			for i := j; i < 4; i++ {
				lk.dataPackage[5+(i*3)+0] = lk.colors[i%j].Red
				lk.dataPackage[5+(i*3)+1] = lk.colors[i%j].Green
				lk.dataPackage[5+(i*3)+2] = lk.colors[i%j].Blue
			}
		}
	}

	return nil
}

func (lk *lKeyboard) setWaveDir() {
	if lk.effect == EffectWave {
		switch lk.waveRTL {
		case RTL:
			lk.dataPackage[18] = 0x01
			lk.dataPackage[19] = 0x00
		case LTR:
			lk.dataPackage[18] = 0x00
			lk.dataPackage[19] = 0x01
		}
	} else {
		lk.dataPackage[18] = 0x00
		lk.dataPackage[19] = 0x00
	}
}

func (lk *lKeyboard) setEmptyBufferBytes() {
	// EMPTY BUFFER
	lk.dataPackage[17] = 0x00
	for i := 20; i < len(lk.dataPackage); i++ {
		lk.dataPackage[i] = 0x00
	}
}

func (lk *lKeyboard) build() error {
	lk.dataPackage = make([]byte, 32)
	// Set Header
	lk.header()
	// EFFECT : type, speed brightness
	lk.setEffect()
	// color
	if err := lk.setColor(); err != nil {
		return err
	}
	lk.setWaveDir()
	lk.setEmptyBufferBytes()
	return nil
}
