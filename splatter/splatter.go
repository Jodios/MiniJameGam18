package splatter

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/utils"
	"math"
	"math/rand"
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
	audioContext *audio.Context
	splatSounds  []*audio.Player
	points       []*Splat
	sprites      map[string]utils.ImageWithFrameDetails
	counter      int
	normalSpeed  int
	backupSpeed  int
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
		sprites:      sprites,
		audioContext: audioContext,
		splatSounds:  contexts,
		normalSpeed:  3,
		backupSpeed:  3,
		points:       make([]*Splat, 0),
	}
}

func (s *Splatter) Update() error {
	s.counter = (s.counter + 1) % math.MaxInt
	return nil
}

func (s *Splatter) Draw(screen *ebiten.Image) {

	randomNumber := rand.Intn(1000)

	// TODO: this is just some bollocks i made up.
	// need to come up with a better way of spawning in something
	// based on speed :/
	spawnCheck := randomNumber > 50 && randomNumber < 53 ||
		randomNumber > 60 && randomNumber < 63 ||
		randomNumber > 0 && randomNumber < 10

	if spawnCheck && !s.splatting {
		s.points = append(s.points, createRandomSplat())
	}

	for _, splat := range s.points {
		s.splatting = false
		opts := &ebiten.DrawImageOptions{}
		frame := s.sprites[splat.asset]
		frameWidth, frameHeight := float64(frame.FrameData.SourceSize.W)*splat.scale, float64(frame.FrameData.SourceSize.H)*splat.scale

		opts.GeoM.Scale(splat.scale, splat.scale)
		opts.GeoM.Translate(
			splat.X-frameWidth/2,
			splat.Y-frameHeight/2,
		)
		screen.DrawImage(frame.Image, opts)
		if splat.done {
			continue
		}
		spawnSpeed := 0.09
		if splat.scale > splat.normalScale && splat.scale-spawnSpeed >= splat.normalScale {
			splat.scale -= spawnSpeed
			s.splatting = true
		} else {
			randomSound := rand.Intn(len(s.splatSounds))
			s.splatSounds[randomSound].Rewind()
			s.splatSounds[randomSound].Play()
			splat.done = true
		}
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
