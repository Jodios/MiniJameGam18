package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jodios/minijamegame18/assets/sprites"
	"github.com/jodios/minijamegame18/brushes"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/screens"
	"github.com/jodios/minijamegame18/utils"
	"log"
)

type STATE int

const (
	START STATE = iota
	SWEEP
)

type Game struct {
	audioContext       *audio.Context
	sprites            map[string]utils.ImageWithFrameDetails
	state              STATE
	startButtonClicked bool
	startButtonHover   bool
	counter            int
	startScreen        *screens.StartScreen
	map1               *screens.Level
	brush              *brushes.Brush
}

func (g *Game) Update() error {
	switch g.state {
	case START:
		g.startScreen.Update()
	case SWEEP:
		// ideally calculations should happen in update function
		g.map1.Update()
		g.brush.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case START:
		g.startScreen.Draw(screen)
		if g.startScreen.DONE {
			g.state = SWEEP
		}
	case SWEEP:
		g.map1.Draw(screen)
		g.brush.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(windowWidth, windowHeight int) (resWidth, resHeight int) {
	return constants.ResX, constants.ResY
}

func main() {
	ebiten.SetWindowSize(constants.ResX*4, constants.ResY*4)
	ebiten.SetWindowTitle("Mop It Up!")

	// setting up audio context
	audioContext := audio.NewContext(constants.AudioSampleRate)

	// unpacking sprites packed by texture packer
	unpacker := &utils.Unpacker{}
	sprites, err := unpacker.UnpackWithFrameDetails(sprites.SPRITES_CONFIG, sprites.SPRITES_PNG)
	check(err)

	game := &Game{
		map1:         screens.NewLevelScreen(audioContext, sprites),
		audioContext: audioContext,
		sprites:      sprites,
		startScreen:  screens.NewStartScreen(audioContext, sprites),
		brush:        brushes.NewBrush(audioContext, sprites),
		//state:        SWEEP,
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
