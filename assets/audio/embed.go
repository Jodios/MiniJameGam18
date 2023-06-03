package audio

import _ "embed"

var (
	//go:embed gotta_sweep.mp3
	GOTTA_SWEEP []byte

	//go:embed sweepin_time.mp3
	SWEEPIN_TIME []byte

	//go:embed song_01.mp3
	INTRO_SONG []byte

	//go:embed song_02_spedup.mp3
	MAIN_LOOP_SONG []byte
)
