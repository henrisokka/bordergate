package main

import "github.com/hajimehoshi/ebiten/v2"

type object struct {
	coord        coord
	spriteOffset coord
	sprite       *ebiten.Image
	handler      func(*Game)
}

func dialogHandlerFactory(dialogID string) func(*Game) {
	return func(g *Game) {
		g.dialogChain = g.dialogs[dialogID]
	}
}
