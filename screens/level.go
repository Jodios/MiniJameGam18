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
	"github.com/jodios/minijamegame18/public"
	"github.com/jodios/minijamegame18/splatter"
	"github.com/jodios/minijamegame18/utils"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Level struct {
	sprites        map[string]utils.ImageWithFrameDetails
	background     *utils.Map
	audioContext   *audio.Context
	gameplaySong   *audio.Player
	counter        int
	splatGenerator *splatter.Splatter
	score          int
	health         int
	Font           font.Face
	DONE           bool
}

func NewLevelScreen(audioContext *audio.Context, sprites map[string]utils.ImageWithFrameDetails) *Level {
	mainSongDecoded, err := mp3.DecodeWithoutResampling(bytes.NewReader(audio2.MAIN_LOOP_SONG))
	check(err)

	mainSong, err := audioContext.NewPlayer(mainSongDecoded)
	check(err)

	sc, err := utils.GetMapConfig(tiles.MAP_1_CONFIG, tiles.MAP_1)
	check(err)

	mf, err := opentype.Parse(public.MapleMono)
	check(err)
	font, err := opentype.NewFace(mf, &opentype.FaceOptions{
		Size:    20,
		DPI:     80,
		Hinting: font.HintingVertical,
	})
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
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && s.CheckCollision(ebiten.CursorPosition()) {
			l.score++
			l.splatGenerator.Speed = l.score/10 + 1
			l.removeSplat(i)
		}
	}
	return nil
}

func (l *Level) removeSplat(i int) {
	l.splatGenerator.Splats =
		append(l.splatGenerator.Splats[:i], l.splatGenerator.Splats[i+1:]...)
}

func (l *Level) Draw(screen *ebiten.Image) {
	l.background.Draw(screen)
	l.splatGenerator.Draw(screen)
	scoreString := fmt.Sprintf("Score: %d", l.score)
	healthString := fmt.Sprintf("Air Quality: %d", l.health)
	text.Draw(screen, scoreString, l.Font, 0, 20, color.Black)
	text.Draw(screen, healthString, l.Font, 0, 40, color.Black)
}
