package audio

import _ "embed"

var (
	//go:embed gotta_sweep.mp3
	GOTTA_SWEEP []byte

	//go:embed sweepin_time.mp3
	SWEEPIN_TIME []byte
)
