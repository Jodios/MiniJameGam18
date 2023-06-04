package screens

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	audio2 "github.com/jodios/minijamegame18/assets/audio"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/splatter"
	"github.com/jodios/minijamegame18/utils"
	"golang.org/x/image/font"
)

type Level struct {
	Score          int
	Font           font.Face
	DONE           bool
	sprites        map[string]utils.ImageWithFrameDetails
	background     *utils.Map
	audioContext   *audio.Context
	gameplaySong   *audio.Player
	counter        int
	splatGenerator *splatter.Splatter
	health         int
}

func NewLevelScreen(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails, font font.Face) *Level {
	mainSongDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.MAIN_LOOP_SONG))
	check(err)

	mainSong, err := audioContext.NewPlayer(mainSongDecoded)
	check(err)

	sc, err := utils.GetMapConfig(tiles.MAP_1_CONFIG, tiles.MAP_1)
	check(err)

	return &Level{
		audioContext:   audioContext,
		background:     sc,
		sprites:        sprites,
		gameplaySong:   mainSong,
		health:         5,
		splatGenerator: splatter.NewSplatter(audioContext, sprites),
		Font:           font,
	}
}

func (l *Level) Update() error {
	if l.health <= 0 {
		l.DONE = true
		l.gameplaySong.Close()
		return nil
	}
	l.counter = (l.counter + 1) % math.MaxInt
	if !l.gameplaySong.IsPlaying() {
		l.gameplaySong.Rewind()
		l.gameplaySong.Play()
	}
	l.splatGenerator.Update()
	for i, s := range l.splatGenerator.Splats {
		splatIsTooOld := time.Now().Sub(s.CreationTime) > l.splatGenerator.SplatLifetime
		if splatIsTooOld && !s.IsMoldy {
			s.MakeMoldy()
			l.health--
		}
		if (inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) || inpututil.IsTouchJustReleased(ebiten.TouchID(0))) &&
			(s.CheckCollision(ebiten.CursorPosition()) || s.CheckCollision(inpututil.TouchPositionInPreviousTick(ebiten.TouchID(0)))) {

			l.Score++
			l.splatGenerator.Speed = l.Score/10 + 1
			l.removeSplat(i)

		}
	}
	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	l.background.Draw(screen)
	l.splatGenerator.Draw(screen)
	scoreString := fmt.Sprintf("Score: %d", l.Score)
	healthString := fmt.Sprintf("Air Quality: %d", l.health)
	text.Draw(screen, scoreString, l.Font, 0, 20, color.Black)
	text.Draw(screen, healthString, l.Font, 0, 40, color.Black)
}

func (l *Level) removeSplat(i int) {
	// when clicking on 2 things at a time it crashes
	// so validating the index to avoid this :)
	if i < 0 || i >= len(l.splatGenerator.Splats) {
		return
	}
	l.splatGenerator.Splats =
		append(l.splatGenerator.Splats[:i], l.splatGenerator.Splats[i+1:]...)
}
