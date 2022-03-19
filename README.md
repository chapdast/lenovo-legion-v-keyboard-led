# Lenovo Legion 5 Keyboard LED Controler

This simple app is just go implementation of a python script by [InstinctEx](https://github.com/InstinctEx/lenovo-ideapad-legion-keyboard-led), will let you control LED light effects of Lenovo Legion 5 Keyboard.
This uses go implementation of libusb-1.0 form
[gousb](https://github.com/google/gousb)



## Effects Avaliable are:
 static
 breath
 hue
 wave

## Install:
clone this repository and build it.

``` go build -o build/lvl cmd/main.go```

## Usage:
Complete command is:

```
lvl -e [effect] -s [speed] -b [brightness] -c [colors-list-by-comma] -d [wave-effect-direction]
```
but some effect may use some of options for example
```
lvl -e static -c 00ff00,0000ff
```
will set keyboard light effect a static effect with two color of Green and Blue repeated across all 4 sections.



You can use -h option to see below help message.
```
Usage of lvl:
  -b int
    	brightness of LEDs
    		0-off
    		1
    		2-brightest (default 2)
  -bright int
    	brightness of LEDs
    		0-off
    		1
    		2-brightest (default 2)
  -c string
    	effect color
    		eg.: 00ff00 | ff00ff,00fff
  -colors string
    	effect color
    		eg.: 00ff00 | ff00ff,00fff
  -d string
    	wave effect direction
    		rtl
    		ltr (default "ltr")
  -debug
    	show debug messages
  -dir string
    	wave effect direction
    		rtl
    		ltr (default "ltr")
  -e string
    	Keyboard LED light effect:
    		static 
    		breath
    		wave
    		hue
    		off (default "off")
  -effect string
    	Keyboard LED light effect:
    		static 
    		breath
    		wave
    		hue
    		off (default "off")
  -s int
    	effect speed
    		1-slowest
    		2
    		3
    		4-fastest (default 1)
  -speed int
    	effect speed
    		1-slowest
    		2
    		3
    		4-fastest (default 1)
  -vpid string
    	Vendor and Product ID in VID:PID format. (default "048d:c955")

```