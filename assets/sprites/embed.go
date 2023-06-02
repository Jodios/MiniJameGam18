package sprites

import _ "embed"

var (
	//go:embed sprites.json
	SPRITES_CONFIG []byte

	//go:embed sprites.png
	SPRITES_PNG []byte
)
