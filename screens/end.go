package screens

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	color2 "image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/utils"
)

type EndScreen struct {
	sprites         map[string]utils.ImageWithFrameDetails
	background      *utils.Map
	audioContext    *audio.Context
	endSong         *audio.Player
	sweepinTime     *audio.Player
	endButton       utils.ImageWithFrameDetails
	endButtonHover  bool
	endButtonX      float64
	endButtonY      float64
	endButtonWidth  float64
	endButtonHeight float64
	counter         int
	font            font.Face
	Score           int
	DONE            bool
}

func NewEndScreen(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails, font font.Face) *EndScreen {
	introSongDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.INTRO_SONG))
	check(err)

	introSong, err := audioContext.NewPlayer(introSongDecoded)
	check(err)

	sweepinTimeDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.SWEEPIN_TIME))
	check(err)

	sweepinTime, err := audioContext.NewPlayer(sweepinTimeDecoded)
	check(err)

	// I am using the same background as start screen; hence the file name
	sc, err := utils.GetMapConfig(tiles.START_SCREEN_1_CONFIG, tiles.START_SCREEN_1_PNG)
	check(err)

	endButton := sprites["restart_button.png"]
	endButtonWidth := float64(endButton.FrameData.SourceSize.W)
	endButtonHeight := float64(endButton.FrameData.SourceSize.H)
	return &EndScreen{
		audioContext:    audioContext,
		background:      sc,
		sprites:         sprites,
		font:            font,
		endSong:         introSong,
		sweepinTime:     sweepinTime,
		endButtonX:      constants.ResX/2 - endButtonWidth/2,
		endButtonY:      constants.ResY/2 + endButtonHeight/2,
		endButtonWidth:  endButtonWidth,
		endButtonHeight: endButtonHeight,
		endButton:       endButton,
	}
}

func (e *EndScreen) Update() error {
	e.counter = (e.counter + 1) % math.MaxInt
	if !e.endSong.IsPlaying() && !e.DONE {
		e.endSong.Rewind()
		e.endSong.Play()
	}
	// checking if mouse is hovering over end button
	mouseX, mouseY := ebiten.CursorPosition()
	tapX, tapY := inpututil.TouchPositionInPreviousTick(ebiten.TouchID(0))
	mouseHover := float64(mouseX) > e.endButtonX &&
		float64(mouseX) < e.endButtonX+e.endButtonWidth &&
		float64(mouseY) > e.endButtonY &&
		float64(mouseY) < e.endButtonY+e.endButtonHeight
	tapHover := float64(tapX) > e.endButtonX &&
		float64(tapX) < e.endButtonX+e.endButtonWidth &&
		float64(tapY) > e.endButtonY &&
		float64(tapY) < e.endButtonY+e.endButtonHeight
	e.endButtonHover = mouseHover || tapHover

	if e.endButtonHover {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}

	if (inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) || inpututil.IsTouchJustReleased(ebiten.TouchID(0))) && e.endButtonHover {
		e.sweepinTime.Rewind()
		e.sweepinTime.Play()
		e.endSong.Close()
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
		e.DONE = true
	}
	return nil
}

func (e *EndScreen) Draw(screen *ebiten.Image) {
	e.background.Draw(screen)

	yOffset := math.Sin(float64((e.counter / 5) % 100))
	opts := new(ebiten.DrawImageOptions)
	//opts.Filter = ebiten.FilterLinear
	opts.GeoM.Translate(e.endButtonX, e.endButtonY+yOffset)
	if e.endButtonHover {
		opts.GeoM.Translate(3, 3)
	}
	screen.DrawImage(e.endButton.Image, opts)

	color := color2.RGBA{
		R: 0x51,
		G: 0x4b,
		B: 0x6d,
		A: 0xff,
	}
	messageStartX := constants.ResX / 2
	messageStartY := constants.ResY / 2
	m1 := fmt.Sprintf(fmt.Sprintf("You've been"))
	m2 := fmt.Sprintf("poisoned by mold!")
	m3 := fmt.Sprintf("Score: %d", e.Score)
	text.Draw(screen, m1, e.font, messageStartX/2-10, messageStartY/3, color)
	text.Draw(screen, m2, e.font, messageStartX/2-48, messageStartY/3+30, color)
	text.Draw(screen, m3, e.font, messageStartX/2+10, messageStartY/3+70, color)
}
