package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/assets/sprites"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/utils"
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
	counter            int
	brushSpeed         int
	brushScale         int
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	mouseX, mouseY := ebiten.CursorPosition()
	switch g.state {
	case START:
		startButton := g.sprites["start_button.png"]
		startButtonWidth := float64(startButton.FrameData.SourceSize.W)
		startButtonHeight := float64(startButton.FrameData.SourceSize.H)

		buttonPositionX, buttonPositionY := float64(ResX/2-startButtonWidth/2), float64(ResY/2+startButtonHeight/2)
		opts.GeoM.Translate(buttonPositionX, buttonPositionY)

		// checking if mouse is hovering over start button
		mouseIsHoveringOverStart := float64(mouseX) > buttonPositionX &&
			float64(mouseX) < buttonPositionX+startButtonWidth &&
			float64(mouseY) > buttonPositionY &&
			float64(mouseY) < buttonPositionY+startButtonHeight

		if mouseIsHoveringOverStart {
			g.startButtonHover = true
			opts.GeoM.Translate(3, 3)
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
		} else {
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
			g.startButtonHover = false
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && mouseIsHoveringOverStart {
			g.startButtonClicked = true
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && g.startButtonClicked {
			g.gottaSweep.Rewind()
			g.gottaSweep.Play()
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
			g.state = SWEEP
		}
		screen.DrawImage(startButton.Image, opts)
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
		brushSpeed:   6,
		brushScale:   3,
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
