package screens

import (
	"bytes"
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

type StartScreen struct {
	sprites           map[string]utils.ImageWithFrameDetails
	background        *utils.Map
	audioContext      *audio.Context
	introSong         *audio.Player
	sweepinTime       *audio.Player
	startButton       utils.ImageWithFrameDetails
	startButtonHover  bool
	startButtonX      float64
	startButtonY      float64
	startButtonWidth  float64
	startButtonHeight float64
	counter           int
	DONE              bool
}

func NewStartScreen(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails) *StartScreen {
	introSongDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.INTRO_SONG))
	check(err)

	introSong, err := audioContext.NewPlayer(introSongDecoded)
	check(err)

	sweepinTimeDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.SWEEPIN_TIME))
	check(err)

	sweepinTime, err := audioContext.NewPlayer(sweepinTimeDecoded)
	check(err)

	sc, err := utils.GetMapConfig(tiles.START_SCREEN_1_CONFIG, tiles.START_SCREEN_1_PNG)
	check(err)

	startButton := sprites["start_button.png"]
	startButtonWidth := float64(startButton.FrameData.SourceSize.W)
	startButtonHeight := float64(startButton.FrameData.SourceSize.H)
	return &StartScreen{
		audioContext:      audioContext,
		background:        sc,
		sprites:           sprites,
		introSong:         introSong,
		sweepinTime:       sweepinTime,
		startButtonX:      constants.ResX/2 - startButtonWidth/2,
		startButtonY:      constants.ResY/2 + startButtonHeight/2,
		startButtonWidth:  startButtonWidth,
		startButtonHeight: startButtonHeight,
		startButton:       startButton,
	}
}

func (s *StartScreen) Update() error {
	s.counter = (s.counter + 1) % math.MaxInt
	if !s.introSong.IsPlaying() && !s.DONE {
		s.introSong.Rewind()
		s.introSong.Play()
	}

	// checking if mouse/touch is hovering over start button
	mouseX, mouseY := ebiten.CursorPosition()
	tapX, tapY := inpututil.TouchPositionInPreviousTick(ebiten.TouchID(0))
	mouseHover := float64(mouseX) > s.startButtonX &&
		float64(mouseX) < s.startButtonX+s.startButtonWidth &&
		float64(mouseY) > s.startButtonY &&
		float64(mouseY) < s.startButtonY+s.startButtonHeight
	tapHover := float64(tapX) > s.startButtonX &&
		float64(tapX) < s.startButtonX+s.startButtonWidth &&
		float64(tapY) > s.startButtonY &&
		float64(tapY) < s.startButtonY+s.startButtonHeight
	s.startButtonHover = mouseHover || tapHover

	// checking if tapped on button (touch screen specific)
	if s.startButtonHover {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	if (inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) || inpututil.IsTouchJustReleased(ebiten.TouchID(0))) && s.startButtonHover {
		s.sweepinTime.Rewind()
		s.sweepinTime.Play()
		s.introSong.Close()
		ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
		s.DONE = true
	}
	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	s.background.Draw(screen)
	opts.GeoM.Translate(s.startButtonX, s.startButtonY)
	if s.startButtonHover {
		opts.GeoM.Translate(3, 3)
	}
	screen.DrawImage(s.startButton.Image, opts)

	yOffset := math.Sin(float64((s.counter / 4) % 100))
	title := s.sprites["title.png"]
	titleScaledWidth := float64(title.FrameData.SourceSize.W) * .4
	opts.Filter = ebiten.FilterLinear
	opts.GeoM.Reset()
	opts.GeoM.Scale(.4, .4)
	opts.GeoM.Translate(0, 16)
	opts.GeoM.Translate(constants.ResX/2-titleScaledWidth/2.3, yOffset)
	opts.GeoM.Rotate(15 * (math.Pi / 180))
	screen.DrawImage(title.Image, opts)
}
