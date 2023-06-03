package splatter

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/utils"
)

type Splat struct {
	X           float64
	Y           float64
	asset       string
	normalScale float64
	scale       float64
	done        bool
}

type Splatter struct {
	Speed        int
	frequency    int
	audioContext *audio.Context
	splatSounds  []*audio.Player
	points       []*Splat
	sprites      map[string]utils.ImageWithFrameDetails
	counter      int
	brush        int
	splatting    bool
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
		Speed:        1,
		sprites:      sprites,
		audioContext: audioContext,
		splatSounds:  contexts,
		frequency:    10,
		points:       make([]*Splat, 0),
	}
}

func (splatManager *Splatter) Update() error {
	splatManager.counter = (splatManager.counter + 1) % math.MaxInt
	frequency := 60 / splatManager.Speed
	spawnCheck := splatManager.counter%frequency == 0

	if spawnCheck && !splatManager.splatting {
		splatManager.points = append(splatManager.points, createRandomSplat())
	}

	for _, splat := range splatManager.points {
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
	for _, splat := range splatManager.points {
		opts := &ebiten.DrawImageOptions{}
		frame := splatManager.sprites[splat.asset]
		frameWidth, frameHeight := float64(frame.FrameData.SourceSize.W)*splat.scale, float64(frame.FrameData.SourceSize.H)*splat.scale

		opts.GeoM.Scale(splat.scale, splat.scale)
		opts.GeoM.Translate(
			splat.X-frameWidth/2,
			splat.Y-frameHeight/2,
		)
		screen.DrawImage(frame.Image, opts)
	}
}

func createRandomSplat() *Splat {
	return &Splat{
		X:           float64(rand.Intn(constants.ResX)),
		Y:           float64(rand.Intn(constants.ResY)),
		asset:       fmt.Sprintf("splat%02d.png", rand.Intn(36)),
		normalScale: .1,
		scale:       1,
	}
}
