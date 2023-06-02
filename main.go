package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/assets/sprites"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/utils"
	"golang.org/x/exp/shiny/materialdesign/colornames"
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
		startButtonWidth := float64(startButton.FrameData.SourceSize.W)
		startButtonHeight := float64(startButton.FrameData.SourceSize.H)

		buttonPositionX, buttonPositionY := float64(ResX/2-startButtonWidth/2), float64(ResY/2+startButtonHeight/2)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(buttonPositionX, buttonPositionY)
		screen.DrawImage(startButton.Image, opts)
		// checking if mouse is hovering over start button
		mouseX, mouseY := ebiten.CursorPosition()

		mouseIsHoveringOverStart := (float64(mouseX) > buttonPositionX &&
			float64(mouseX) < buttonPositionX+startButtonWidth &&
			float64(mouseY) > buttonPositionY &&
			float64(mouseY) < buttonPositionY+startButtonHeight)

		if mouseIsHoveringOverStart {
			g.startButtonHover = true
			fmt.Println("hover", time.Now())
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
