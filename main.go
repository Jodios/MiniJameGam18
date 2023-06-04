package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/jodios/minijamegame18/assets/sprites"
	"github.com/jodios/minijamegame18/brushes"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/public"
	"github.com/jodios/minijamegame18/screens"
	"github.com/jodios/minijamegame18/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

type STATE int

const (
	START STATE = iota
	SWEEP
	END
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
	endScreen          *screens.EndScreen
	font               font.Face
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
	case END:
		g.endScreen.Update()
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
		if g.map1.DONE {
			g.endScreen.Score = g.map1.Score
			g.state = END
		}
	case END:
		g.endScreen.Draw(screen)
		if g.endScreen.DONE {
			g.endScreen = screens.NewEndScreen(g.audioContext, g.sprites, g.font)
			g.map1 = screens.NewLevelScreen(g.audioContext, g.sprites, g.font)
			g.state = SWEEP
		}
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualFPS()))
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

	mf, err := opentype.Parse(public.MapleMono)
	check(err)
	font, err := opentype.NewFace(mf, &opentype.FaceOptions{
		Size:    20,
		DPI:     80,
		Hinting: font.HintingVertical,
	})
	check(err)

	game := &Game{
		map1:         screens.NewLevelScreen(audioContext, sprites, font),
		audioContext: audioContext,
		sprites:      sprites,
		startScreen:  screens.NewStartScreen(audioContext, sprites),
		endScreen:    screens.NewEndScreen(audioContext, sprites, font),
		brush:        brushes.NewBrush(audioContext, sprites),
		font:         font,
		//state:        END,
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
