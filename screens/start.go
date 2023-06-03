package screens

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/constants"
	"github.com/jodios/minijamegame18/utils"
	"log"
)

type StartScreen struct {
	sprites            map[string]utils.ImageWithFrameDetails
	audioContext       *audio.Context
	introSong          *audio.Player
	sweepinTime        *audio.Player
	startButton        utils.ImageWithFrameDetails
	startButtonHover   bool
	startButtonX       float64
	startButtonY       float64
	startButtonWidth   float64
	startButtonHeight  float64
	startButtonClicked bool
	DONE               bool
}

func NewStartScreen(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails) *StartScreen {
	introSongDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.INTRO_SONG))
	if err != nil {
		log.Fatal(err)
	}
	introSong, err := audioContext.NewPlayer(introSongDecoded)
	if err != nil {
		log.Fatal(err)
	}
	sweepinTimeDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.SWEEPIN_TIME))
	if err != nil {
		log.Fatal(err)
	}
	sweepinTime, err := audioContext.NewPlayer(sweepinTimeDecoded)
	if err != nil {
		log.Fatal(err)
	}

	startButton := sprites["start_button.png"]
	startButtonWidth := float64(startButton.FrameData.SourceSize.W)
	startButtonHeight := float64(startButton.FrameData.SourceSize.H)
	return &StartScreen{
		audioContext:      audioContext,
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
	if !s.introSong.IsPlaying() {
		s.introSong.Rewind()
		s.introSong.Play()
	}
	mouseX, mouseY := ebiten.CursorPosition()
	// checking if mouse is hovering over start button
	s.startButtonHover = float64(mouseX) > s.startButtonX &&
		float64(mouseX) < s.startButtonX+s.startButtonWidth &&
		float64(mouseY) > s.startButtonY &&
		float64(mouseY) < s.startButtonY+s.startButtonHeight
	if s.startButtonHover {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && s.startButtonHover {
		s.startButtonClicked = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && s.startButtonClicked {
		s.sweepinTime.Rewind()
		s.sweepinTime.Play()
		s.introSong.Close()
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
		s.DONE = true
	}
	return nil
}
func (s *StartScreen) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(s.startButtonX, s.startButtonY)
	if s.startButtonHover {
		opts.GeoM.Translate(3, 3)
	}
	screen.DrawImage(s.startButton.Image, opts)
}
