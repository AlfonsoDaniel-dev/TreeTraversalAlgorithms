package Game

import (
	"image/color"

	"github.com/AlfonsoDaniel-dev/TreeTraversal/src/Tree"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	BackgroundColor color.Color
	ScreenWidth     int
	ScreenHeight    int
	Tree            *Tree.Tree
}

func NewGame(BackgroundColor color.RGBA, ScreenWidth, ScreenHeight int, tree *Tree.Tree) *Game {
	return &Game{
		BackgroundColor: BackgroundColor,
		ScreenWidth:     ScreenWidth,
		ScreenHeight:    ScreenHeight,
		Tree:            tree,
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.BackgroundColor)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}
