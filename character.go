package main

import "github.com/hajimehoshi/ebiten/v2"

type character struct {
	x             float64
	y             float64
	actionCounter int

	sprites         map[string][]*ebiten.Image
	activeAnimation string

	currentSprite int
}

func (c *character) nextSprite() {
	if c.currentSprite+1 > len(c.sprites[c.activeAnimation])-1 {
		c.currentSprite = 0
	} else {
		c.currentSprite += 1
	}
}
