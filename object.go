package main

import "github.com/hajimehoshi/ebiten/v2"

type object struct {
	coord        coord
	spriteOffset coord
	sprite       *ebiten.Image
	handler      func(*Game)
}
