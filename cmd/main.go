package main

import (
	"image/color"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	ebiten.SetWindowTitle("Graph Test")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	err := ebiten.RunGame(Game.NewGame(color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 1,
	}, 1280, 720, nil))
	if err != nil {
		panic("could not start the game")
	}

}
