package constants

import "image/color"

const (
	ResX            = 256
	ResY            = 256
	AudioSampleRate = 48000
)

var (
	// 0 to 1 !!!! holy shit this blew up my ears 😢😢😢😢
	Volume           = VolumePercentage * 2
	VolumePercentage = .5
	BackgroundColor  = color.RGBA{
		R: 0x51,
		G: 0x4b,
		B: 0x6d,
		A: 0xff,
	}
	SecondaryColor = color.RGBA{
		R: 0xb0,
		G: 0xa7,
		B: 0xc5,
		A: 0xff,
	}
	TertiaryColor = color.RGBA{
		R: 0xe3,
		G: 0xd7,
		B: 0xff,
		A: 0xff,
	}
)
