package tiles

import _ "embed"

var (
	//go:embed Kitchen_and_more_tileset/map1/map1.txt
	MAP_1_CONFIG []byte
	//go:embed Kitchen_and_more_tileset/map1/map1.png
	MAP_1 []byte

	//go:embed Kitchen_and_more_tileset/startScreen/startScreen.txt
	START_SCREEN_1_CONFIG []byte
	//go:embed Kitchen_and_more_tileset/startScreen/startScreen.png
	START_SCREEN_1_PNG []byte
)
