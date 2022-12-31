package main

import "github.com/hajimehoshi/ebiten/v2"

type character struct {
	coord         coord
	actionCounter int

	activeAnimation string
	direction       string
	walking         bool
	hitting         bool

	sprites       map[string][]*ebiten.Image
	currentSprite int
	spriteOffset  coord
}

func (c *character) nextSprite() {
	if !c.walking && !c.hitting {
		return
	}

	if c.currentSprite+1 > len(c.sprites[c.activeAnimation])-1 {
		c.currentSprite = 0
	} else {
		c.currentSprite += 1
	}
}
