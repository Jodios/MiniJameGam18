package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"log"
)

const (
	ResX = 512
	ResY = 512
)

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Blue200)
}

func (g *Game) Layout(windowWidth, windowHeight int) (resWidth, resHeight int) {
	return ResX, ResY
}

func main() {
	ebiten.SetWindowSize(ResX*2, ResY*2)
	ebiten.SetWindowTitle("Mopper")
	game := &Game{}
	log.Fatal(ebiten.RunGame(game))
}
