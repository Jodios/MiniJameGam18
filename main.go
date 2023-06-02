package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/utils"
	"log"
)

const (
	ResX            = 256
	ResY            = 256
	AudioSampleRate = 48000
)

type Game struct {
	map1         *utils.Map
	audioContext *audio.Context
	gottaSweep   *audio.Player
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.map1.Draw(screen)
}

func (g *Game) Layout(windowWidth, windowHeight int) (resWidth, resHeight int) {
	return ResX, ResY
}

func main() {
	ebiten.SetWindowSize(ResX*2, ResY*2)
	ebiten.SetWindowTitle("Mopper")

	// setting up audio context
	audioContext := audio.NewContext(AudioSampleRate)
	gottaSweepDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.SWEEPIN_TIME))
	check(err)
	gottaSweep, err := audioContext.NewPlayer(gottaSweepDecoded)
	check(err)
	gottaSweep.Play()
	//gottaSweep.Rewind()

	m1, err := utils.GetMapConfig(tiles.MAP_1_CONFIG, tiles.MAP_1)
	check(err)
	game := &Game{
		map1:         m1,
		audioContext: audioContext,
		gottaSweep:   gottaSweep,
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
