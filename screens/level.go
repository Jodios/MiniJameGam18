package screens

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/utils"
	"math"
)

type Level struct {
	sprites      map[string]utils.ImageWithFrameDetails
	background   *utils.Map
	audioContext *audio.Context
	gameplaySong *audio.Player
	counter      int
	DONE         bool
}

func NewLevelScreen(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails) *Level {
	mainSongDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.MAIN_LOOP_SONG))
	check(err)

	mainSong, err := audioContext.NewPlayer(mainSongDecoded)
	check(err)

	sc, err := utils.GetMapConfig(tiles.MAP_1_CONFIG, tiles.MAP_1)
	check(err)

	return &Level{
		audioContext: audioContext,
		background:   sc,
		sprites:      sprites,
		gameplaySong: mainSong,
	}
}

func (s *Level) Update() error {
	s.counter = (s.counter + 1) % math.MaxInt
	if !s.gameplaySong.IsPlaying() {
		s.gameplaySong.Rewind()
		s.gameplaySong.Play()
	}
	return nil
}

func (s *Level) Draw(screen *ebiten.Image) {
	s.background.Draw(screen)
}
