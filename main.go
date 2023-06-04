package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jodios/minijamegame18/assets/sprites"
	"github.com/jodios/minijamegame18/brushes"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/public"
	"github.com/jodios/minijamegame18/screens"
	"github.com/jodios/minijamegame18/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"log"
)

type STATE int

const (
	START STATE = iota
	SWEEP
	END
)

type Game struct {
	audioContext        *audio.Context
	sprites             map[string]utils.ImageWithFrameDetails
	state               STATE
	counter             int
	startScreen         *screens.StartScreen
	map1                *screens.Level
	endScreen           *screens.EndScreen
	font                font.Face
	brush               *brushes.Brush
	soundIcon           utils.ImageWithFrameDetails
	soundIconHovered    bool
	soundIconX          float64
	soundIconY          float64
	soundCircleX        float64
	soundCircleY        float64
	soundSettingsX      float64
	soundSettingsY      float64
	soundSettingsWidth  float64
	soundSettingsHeight float64
	soundIconHover      bool
	soundCircleHover    bool
	soundCircleClicked  bool
	settings            bool
	soundCircleRadius   int
}

func (g *Game) Update() error {
	switch g.state {
	case START:
		g.startScreen.Update(g.settings)
	case SWEEP:
		// ideally calculations should happen in update function
		g.map1.Update(g.settings)
		g.brush.Update()
	case END:
		g.endScreen.Update(g.settings)
	}

	if g.soundIconHovered {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	mousePosX, mousePosY := ebiten.CursorPosition()
	if g.settings {
		circleX := g.soundSettingsX + g.soundCircleX
		circleY := g.soundSettingsY + g.soundCircleY
		g.soundCircleHover = float64(mousePosX) > circleX-float64(g.soundCircleRadius) &&
			float64(mousePosX) < circleX+float64(g.soundCircleRadius) &&
			float64(mousePosY) > circleY-float64(g.soundCircleRadius) &&
			float64(mousePosY) < circleY+float64(g.soundCircleRadius)
		if g.soundCircleHover {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
		}
		if (g.soundCircleHover && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)) || g.soundCircleClicked {
			g.soundCircleClicked = true
			newPos := float64(mousePosX) - 32
			if newPos > 16 && newPos < g.soundSettingsWidth-16 {
				g.soundCircleX = newPos
				constants.VolumePercentage = (newPos - 17.0) / 158.0
				constants.Volume = constants.VolumePercentage * 2
			}
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
			g.soundCircleClicked = false
		}
	}
	g.soundIconHover = float64(mousePosX) > g.soundIconX &&
		float64(mousePosX) < g.soundIconX+16 &&
		float64(mousePosY) > g.soundIconY &&
		float64(mousePosY) < g.soundIconY+16

	if g.soundIconHover {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
			g.settings = !g.settings
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case START:
		g.startScreen.Draw(screen, g.settings)
		if g.startScreen.DONE {
			g.state = SWEEP
		}
	case SWEEP:
		g.map1.Draw(screen, g.settings)
		g.brush.Draw(screen)
		if g.map1.DONE {
			g.endScreen.Score = g.map1.Score
			g.state = END
		}
	case END:
		g.endScreen.Draw(screen, g.settings)
		if g.endScreen.DONE {
			g.endScreen = screens.NewEndScreen(g.audioContext, g.sprites, g.font)
			g.map1 = screens.NewLevelScreen(g.audioContext, g.sprites, g.font)
			g.state = SWEEP
		}
	}

	//ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualFPS()))
	opts := new(ebiten.DrawImageOptions)
	opts.GeoM.Translate(g.soundIconX, g.soundIconY)
	screen.DrawImage(g.soundIcon.Image.SubImage(image.Rect(
		0, 0,
		16, g.soundIcon.FrameData.SourceSize.H,
	)).(*ebiten.Image), opts)

	if g.settings {
		soundSettingsImage := g.getSoundThing()
		opts.GeoM.Reset()
		opts.GeoM.Translate(g.soundSettingsX, g.soundSettingsY)
		screen.DrawImage(soundSettingsImage, opts)
	}
}

func (g *Game) Layout(windowWidth, windowHeight int) (resWidth, resHeight int) {
	return constants.ResX, constants.ResY
}

// I really don't feel like making this into its own
// new thing right now lol
func (g *Game) getSoundThing() *ebiten.Image {
	soundSettings := ebiten.NewImage(int(g.soundSettingsWidth), int(g.soundSettingsHeight))
	soundSettings.Fill(constants.BackgroundColor)
	// this is just the bar which will remain constant
	vector.DrawFilledRect(soundSettings, 12, float32(g.soundSettingsHeight/2), float32(g.soundSettingsWidth-24), 4, constants.SecondaryColor, true)

	// circle x range from 16px - (g.soundSettingsWidth - 16)px
	//g.soundCircleX = g.soundSettingsWidth - 16
	vector.DrawFilledCircle(soundSettings, float32(g.soundCircleX), float32(g.soundCircleY), float32(g.soundCircleRadius), constants.TertiaryColor, true)

	opts := new(ebiten.DrawImageOptions)
	opts.ColorScale.ScaleWithColor(constants.TertiaryColor)
	opts.GeoM.Translate(16, 32)
	text.DrawWithOptions(soundSettings, "Volume:", g.font, opts)
	return soundSettings
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
		map1:                screens.NewLevelScreen(audioContext, sprites, font),
		audioContext:        audioContext,
		sprites:             sprites,
		startScreen:         screens.NewStartScreen(audioContext, sprites),
		endScreen:           screens.NewEndScreen(audioContext, sprites, font),
		brush:               brushes.NewBrush(audioContext, sprites),
		font:                font,
		soundIcon:           sprites["sound.png"],
		soundIconX:          10.0,
		soundIconY:          float64(constants.ResY-sprites["sound.png"].FrameData.SourceSize.H) - 10,
		soundSettingsWidth:  constants.ResX - 64,
		soundSettingsHeight: constants.ResY / 3,
		soundCircleX:        (constants.ResX - 64) / 2,
		soundCircleY:        ((constants.ResY / 3) / 2) + 2.8,
		soundSettingsX:      32,
		soundSettingsY:      constants.ResY/2 + 12,
		soundCircleRadius:   5,
		//state:        END,
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
