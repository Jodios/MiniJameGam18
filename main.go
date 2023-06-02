package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/assets/sprites"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/utils"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"log"
)

const (
	ResX            = 256
	ResY            = 256
	AudioSampleRate = 48000
)

type STATE int

const (
	START STATE = iota
	SWEEP
)

type Game struct {
	map1               *utils.Map
	audioContext       *audio.Context
	gottaSweep         *audio.Player
	sprites            map[string]utils.ImageWithFrameDetails
	state              STATE
	startButtonClicked bool
	startButtonHover   bool
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.White)
	switch g.state {
	case START:
		startButton := g.sprites["start_button.png"]
		dx, dy := float64(ResX/2-startButton.FrameData.SourceSize.W/2), float64(ResY/2+startButton.FrameData.SourceSize.H/2)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(dx, dy)
		screen.DrawImage(startButton.Image, opts)
		// checking if mouse is hovering over start button
		mx, my := ebiten.CursorPosition()
		if float64(mx) > dx && mx < mx+startButton.FrameData.SourceSize.W && float64(my) > dy && my < my+startButton.FrameData.SourceSize.H {
			g.startButtonHover = true
			fmt.Println("hover")
		} else {
			g.startButtonHover = false
		}
	case SWEEP:
		g.map1.Draw(screen)
	}
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

	// unpacking sprites packed by texture packer
	unpacker := &utils.Unpacker{}
	sprites, err := unpacker.UnpackWithFrameDetails(sprites.SPRITES_CONFIG, sprites.SPRITES_PNG)
	check(err)

	m1, err := utils.GetMapConfig(tiles.MAP_1_CONFIG, tiles.MAP_1)
	check(err)
	game := &Game{
		map1:         m1,
		audioContext: audioContext,
		gottaSweep:   gottaSweep,
		sprites:      sprites,
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
