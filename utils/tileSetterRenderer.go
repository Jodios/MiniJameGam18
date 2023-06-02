package utils

import (
	"bytes"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
)

type TileSetterJSON struct {
	TileSize  int                `json:"tile_size,omitempty"`
	MapWidth  int                `json:"map_width,omitempty"`
	MapHeight int                `json:"map_height,omitempty"`
	Layers    []*TileSetterLayer `json:"layers,omitempty"`
}
type TileSetterLayer struct {
	Name  string `json:"name,omitempty"`
	Tiles []Tile `json:"positions,omitempty"`
}
type Tile struct {
	X  int `json:"x,omitempty"`
	Y  int `json:"y,omitempty"`
	Id int `json:"id,omitempty"`
}
type Map struct {
	Config *TileSetterJSON
	Image  *ebiten.Image
}

func GetMapConfig(rawConfig, rawAssets []byte) (*Map, error) {
	config := &TileSetterJSON{}
	err := json.Unmarshal(rawConfig, config)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(rawAssets))
	if err != nil {
		return nil, err
	}
	return &Map{Image: ebiten.NewImageFromImage(img), Config: config}, nil
}

func (m *Map) Draw(screen *ebiten.Image) {
	for i := 0; i < len(m.Config.Layers); i++ {
		m.draw(screen, i)
	}
}

func (m *Map) DrawLayersN(screen *ebiten.Image, n ...int) {
	for _, i := range n {
		if i >= len(m.Config.Layers) {
			continue
		}
		m.draw(screen, i)
	}
}

func (m *Map) draw(screen *ebiten.Image, i int) {
	for _, tile := range m.Config.Layers[i].Tiles {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(tile.X*m.Config.TileSize), float64(tile.Y*m.Config.TileSize))
		dx := (tile.Id * m.Config.TileSize) % m.Image.Bounds().Dx()
		dy := (tile.Id * m.Config.TileSize) / m.Image.Bounds().Dx()
		screen.DrawImage(m.Image.SubImage(image.Rect(
			dx, dy*m.Config.TileSize,
			dx+m.Config.TileSize, dy*m.Config.TileSize+m.Config.TileSize,
		)).(*ebiten.Image), opts)
	}
}
