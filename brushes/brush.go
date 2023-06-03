package brushes

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	normalScale     int
	backupScale     int
	normalSpeed     int
	backupSpeed     int
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
		normalSpeed:     3,
		backupSpeed:     3,
		normalScale:     5,
		backupScale:     5,
		sweepSweepSweep: gottaSweep,
	}
}

func (b *Brush) Update() error {
	b.counter = (b.counter + 1) % math.MaxInt
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		b.backupScale = b.normalScale
		b.normalScale = b.normalScale / 2
		b.backupSpeed = b.normalSpeed
		b.normalSpeed = 1
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		b.normalScale = b.backupScale
		b.normalSpeed = b.backupSpeed
	}
	return nil
}

func (b *Brush) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	frameWidth := 16
	scaledWidth := frameWidth * b.normalScale
	i := (b.counter / b.normalSpeed) % (b.sprites[Main].FrameData.SourceSize.W / frameWidth)

	mouseX, mouseY := ebiten.CursorPosition()
	opts.GeoM.Scale(float64(b.normalScale), float64(b.normalScale))
	opts.GeoM.Translate(float64(mouseX+scaledWidth/4), float64(mouseY+scaledWidth/10))
	opts.GeoM.Translate(float64(-b.normalScale*frameWidth), float64(-b.normalScale*frameWidth))
	screen.DrawImage(b.sprites[Main].Image.SubImage(image.Rect(
		i*frameWidth, 0,
		i*frameWidth+frameWidth, frameWidth,
	)).(*ebiten.Image), opts)
}
