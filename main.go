package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jodios/minijamegame18/assets/tiles"
	"github.com/jodios/minijamegame18/utils"
	"log"
)

const (
	ResX = 256
	ResY = 256
)

type Game struct {
	map1 *utils.Map
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.map1.Draw(screen)
}

func (g *Game) Layout(windowWidth, windowHeight int) (resWidth, resHeight int) {
	return ResX, ResY
}

func main() {
	ebiten.SetWindowSize(ResX*2, ResY*2)
	ebiten.SetWindowTitle("Mopper")
	m1, err := utils.GetMapConfig(tiles.MAP_1_CONFIG, tiles.MAP_1)
	check(err)
	game := &Game{
		map1: m1,
	}
	log.Fatal(ebiten.RunGame(game))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
