package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type npc struct {
	activeAnimation string
	currentSprite   int
	sprites         map[string][]*ebiten.Image

	coord coord
	x     int
	y     int

	decisionCounter int
	game            *Game
}

func (n *npc) init(game *Game) {
	n.game = game
	n.decisionCounter = 20

	npcSprites, _, err := ebitenutil.NewImageFromFile("assets/npc.png")
	if err != nil {
		log.Fatal(err)
	}

	sprites := make(map[string][]*ebiten.Image)
	sprites["down"] = splitSprites(npcSprites, 0, 0, 16, 32, 4)
	sprites["right"] = splitSprites(npcSprites, 0, 32, 16, 32, 4)
	sprites["up"] = splitSprites(npcSprites, 0, 64, 16, 32, 4)
	sprites["left"] = splitSprites(npcSprites, 0, 96, 16, 32, 4)
	n.sprites = sprites

	fmt.Println(sprites)
}
func (n *npc) update() {
	if n.activeAnimation == "down" {
		n.coord.y += 1
	}

	if n.activeAnimation == "up" {
		n.coord.y -= 1
	}

	if n.activeAnimation == "right" {
		n.coord.x += 1
	}

	if n.activeAnimation == "left" {
		n.coord.x -= 1
	}

	if n.coord.x < 0 {
		n.decisionCounter = 120
		n.activeAnimation = "right"
	}

	if n.coord.y < 0 {
		n.decisionCounter = 120
		n.activeAnimation = "down"
	}

	if n.decisionCounter < 1 {
		n.decisionCounter = 60

		// TODO: Make decision here

		rand.Seed(time.Now().UnixNano())

		switch rand.Intn(4) {
		case 0:
			n.activeAnimation = "down"
		case 1:
			n.activeAnimation = "up"
		case 2:
			n.activeAnimation = "left"
		case 3:
			n.activeAnimation = "right"
		}
	}

	n.decisionCounter -= 1
}

func (n *npc) nextSprite() {
	if n.currentSprite+1 > len(n.sprites[n.activeAnimation])-1 {
		n.currentSprite = 0
	} else {
		n.currentSprite += 1
	}
}
