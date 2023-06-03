package main

import (
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/screens"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/jodios/minijamegame18/assets/sprites"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/utils"
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
	counter            int
	brushSpeed         int
	brushScale         int
	startScreen        *screens.StartScreen
}

func (g *Game) Update() error {
	switch g.state {
	case START:
		g.startScreen.Update()
	case SWEEP:
		// ideally calculations should happen in update function
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	mouseX, mouseY := ebiten.CursorPosition()
	switch g.state {
	case START:
		g.startScreen.Draw(screen)
		if g.startScreen.DONE {
			g.state = SWEEP
		}
	case SWEEP:
		brush := g.sprites["mops.png"]
		g.counter = (g.counter + 1) % math.MaxInt
		g.map1.Draw(screen)

		frameWidth := 16
		i := (g.counter / g.brushSpeed) % (brush.FrameData.SourceSize.W / frameWidth)

		opts.GeoM.Scale(float64(g.brushScale), float64(g.brushScale))
		opts.GeoM.Translate(float64(-g.brushScale*frameWidth), float64(-g.brushScale*frameWidth))
		opts.GeoM.Translate(float64(mouseX+frameWidth), float64(mouseY+frameWidth))
		screen.DrawImage(brush.Image.SubImage(image.Rect(
			i*frameWidth, 0,
			i*frameWidth+frameWidth, frameWidth,
		)).(*ebiten.Image), opts)
	}
}

func (g *Game) Layout(windowWidth, windowHeight int) (resWidth, resHeight int) {
	return constants.ResX, constants.ResY
}

func main() {
	ebiten.SetWindowSize(constants.ResX*2, constants.ResY*2)
	ebiten.SetWindowTitle("Mop It Up!")

	// setting up audio context
	audioContext := audio.NewContext(constants.AudioSampleRate)

	// unpacking sprites packed by texture packer
	unpacker := &utils.Unpacker{}
	sprites, err := unpacker.UnpackWithFrameDetails(sprites.SPRITES_CONFIG, sprites.SPRITES_PNG)
	check(err)

	m1, err := utils.GetMapConfig(tiles.MAP_1_CONFIG, tiles.MAP_1)
	check(err)
	game := &Game{
		map1:         m1,
		audioContext: audioContext,
		sprites:      sprites,
		brushSpeed:   6,
		brushScale:   3,
		startScreen:  screens.NewStartScreen(audioContext, sprites),
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
