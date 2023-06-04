package splatter

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/utils"
)

var (
	colors = []color.RGBA{
		{0x9a, 0x38, 0x2d, 0xff}, // ketchup red
		{0xe1, 0xad, 0x01, 0xff}, // mustard yellow dark
		{0xff, 0xdb, 0x58, 0xff}, // mustard yellow light
		{0xb3, 0xcb, 0xaa, 0xff}, // relish green light
		{0x6a, 0xa8, 0x4f, 0xff}, // relish green dark
		{0x7c, 0x53, 0x4d, 0xff}, // barbecue sawce
		{0x28, 0x12, 0x1c, 0xff}, // tamarind
		{0x55, 0x1a, 0x8b, 0xff}, // purple ketchup
	}
)

type Splat struct {
	x            float64
	y            float64
	asset        utils.ImageWithFrameDetails
	CreationTime time.Time
	normalScale  float64
	scale        float64
	Color        color.RGBA
	done         bool
	IsMoldy      bool
}

type Splatter struct {
	Speed         int
	frequency     int
	audioContext  *audio.Context
	splatSounds   []*audio.Player
	Splats        []*Splat
	SplatLifetime time.Duration
	sprites       map[string]utils.ImageWithFrameDetails
	counter       int
	brush         int
	splatting     bool
}

func NewSplatter(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails) *Splatter {
	contexts := make([]*audio.Player, 0)
	splatDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.SPLAT_01))
	check(err)
	splatSound, err := audioContext.NewPlayer(splatDecoded)
	check(err)
	contexts = append(contexts, splatSound)
	splatDecoded, err = mp3.DecodeWithoutResampling(bytes.NewReader(audio2.SPLAT_01))
	check(err)
	splatSound, err = audioContext.NewPlayer(splatDecoded)
	check(err)
	contexts = append(contexts, splatSound)
	return &Splatter{
		Speed:         1,
		sprites:       sprites,
		audioContext:  audioContext,
		splatSounds:   contexts,
		SplatLifetime: time.Millisecond * 3000,
		frequency:     1,
		Splats:        make([]*Splat, 0),
	}
}

func (splatManager *Splatter) Update() error {
	splatManager.counter = (splatManager.counter + 1) % math.MaxInt
	frequency := 120 / splatManager.Speed
	spawnCheck := splatManager.counter%frequency == 0

	if spawnCheck && !splatManager.splatting {
		splatManager.Splats = append(splatManager.Splats, splatManager.createRandomSplat())
	}

	for _, splat := range splatManager.Splats {
		splatManager.splatting = false
		if splat.done {
			continue
		}
		spawnSpeed := 0.09
		if splat.scale > splat.normalScale && splat.scale-spawnSpeed >= splat.normalScale {
			splat.scale -= spawnSpeed
			splatManager.splatting = true
		} else {
			randomSound := rand.Intn(len(splatManager.splatSounds))
			splatManager.splatSounds[randomSound].Rewind()
			splatManager.splatSounds[randomSound].Play()
			splat.done = true
		}
	}

	return nil
}

func (splatManager *Splatter) Draw(screen *ebiten.Image) {
	for _, splat := range splatManager.Splats {
		opts := &colorm.DrawImageOptions{}
		frame := splat.asset
		frameWidth, frameHeight := float64(frame.FrameData.SourceSize.W)*splat.scale,
			float64(frame.FrameData.SourceSize.H)*splat.scale

		opts.GeoM.Scale(splat.scale, splat.scale)
		opts.GeoM.Translate(
			splat.x-frameWidth/2,
			splat.y-frameHeight/2,
		)

		colorManager := colorm.ColorM{}
		colorManager.Scale(0, 0, 0, 1)
		r := float64(splat.Color.R) / 0xff
		g := float64(splat.Color.G) / 0xff
		b := float64(splat.Color.B) / 0xff
		colorManager.Translate(r, g, b, 0)

		colorm.DrawImage(screen, frame.Image, colorManager, opts)
	}
}

func (splatManager *Splatter) createRandomSplat() *Splat {
	random := rand.Intn(len(colors))
	return &Splat{
		x:            float64(rand.Intn(constants.ResX)),
		y:            float64(rand.Intn(constants.ResY)),
		asset:        splatManager.sprites[fmt.Sprintf("splat%02d.png", rand.Intn(36))],
		CreationTime: time.Now(),
		normalScale:  .1,
		scale:        1,
		Color:        colors[random],
	}
}

func (s *Splat) CheckCollision(cursorPosX, cursorPosY int) bool {
	var True bool = false
	if !s.done || s.IsMoldy {
		return True
	}
	frameWidth, frameHeight := float64(s.asset.FrameData.SourceSize.W)*s.scale,
		float64(s.asset.FrameData.SourceSize.H)*s.scale
	splatStartX, splatStartY := s.x-frameWidth/2, s.y-frameHeight/2
	isClickingOnSplat := float64(cursorPosX) > splatStartX && float64(cursorPosX) < splatStartX+frameWidth &&
		float64(cursorPosY) > splatStartY && float64(cursorPosY) < splatStartY+frameHeight
	return isClickingOnSplat
}

func (s *Splat) MakeMoldy() {
	s.IsMoldy = true
	s.Color = color.RGBA{0xff, 0xff, 0xff, 0xff}
}
