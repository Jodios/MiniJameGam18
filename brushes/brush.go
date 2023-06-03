package brushes

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/utils"
	"image"
	"math"
)

const (
	Main string = "mops.png"
)

type Brush struct {
	audioContext    *audio.Context
	sweepSweepSweep *audio.Player
	sprites         map[string]utils.ImageWithFrameDetails
	counter         int
	scale           int
	speed           int
	brush           int
}

func NewBrush(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails) *Brush {
	gottaSweepDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.GOTTA_SWEEP))
	check(err)
	gottaSweep, err := audioContext.NewPlayer(gottaSweepDecoded)
	check(err)
	return &Brush{
		sprites:         sprites,
		audioContext:    audioContext,
		speed:           3,
		scale:           5,
		sweepSweepSweep: gottaSweep,
	}
}

func (b *Brush) Update() error {
	b.counter = (b.counter + 1) % math.MaxInt
	return nil
}

func (b *Brush) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	frameWidth := 16
	i := (b.counter / b.speed) % (b.sprites[Main].FrameData.SourceSize.W / frameWidth)

	mouseX, mouseY := ebiten.CursorPosition()
	opts.GeoM.Scale(float64(b.scale), float64(b.scale))
	opts.GeoM.Translate(float64(-b.scale*frameWidth), float64(-b.scale*frameWidth))
	opts.GeoM.Translate(float64(mouseX+frameWidth), float64(mouseY+frameWidth))
	screen.DrawImage(b.sprites[Main].Image.SubImage(image.Rect(
		i*frameWidth, 0,
		i*frameWidth+frameWidth, frameWidth,
	)).(*ebiten.Image), opts)
}
