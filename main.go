package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var charSprites *ebiten.Image
var actionImg *ebiten.Image

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

type npc struct {
	activeAnimation string
	currentSprite   int
	sprites         map[string][]*ebiten.Image

	x int
	y int

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
		n.y += 1
	}

	if n.activeAnimation == "up" {
		n.y -= 1
	}

	if n.activeAnimation == "right" {
		n.x += 1
	}

	if n.activeAnimation == "left" {
		n.x -= 1
	}

	if n.x < 0 {
		n.decisionCounter = 120
		n.activeAnimation = "right"
	}

	if n.y < 0 {
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
	fmt.Println("nextSprite", n.currentSprite+1 > len(n.sprites[n.activeAnimation])-1)
	if n.currentSprite+1 > len(n.sprites[n.activeAnimation])-1 {
		n.currentSprite = 0
	} else {
		n.currentSprite += 1
	}

	fmt.Println("npc sprite: ", n.currentSprite)
}

func init() {

	var err error
	charSprites, _, err = ebitenutil.NewImageFromFile("assets/character.png")

	if err != nil {
		log.Fatal(err)
	}

	actionImg, _, err = ebitenutil.NewImageFromFile("assets/gopher_action.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	hero           character
	frameCount     int
	terrainSprites map[string]*ebiten.Image

	npcList []*npc
}

const MOVE_DELTA = 1

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.hero.activeAnimation = "left"
		g.hero.x -= MOVE_DELTA
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.hero.activeAnimation = "right"
		g.hero.x += MOVE_DELTA
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.hero.activeAnimation = "down"
		g.hero.y += MOVE_DELTA
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.hero.activeAnimation = "up"
		g.hero.y -= MOVE_DELTA
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		fmt.Println("action is triggered")
		g.hero.actionCounter = 30
		g.hero.activeAnimation = "hit"

	} else if g.hero.actionCounter > 0 {
		g.hero.actionCounter -= 1
	} else {
		g.hero.activeAnimation = "idle"
		g.hero.currentSprite = 0
	}

	if g.hero.actionCounter > 0 {
		g.hero.actionCounter -= 1
	}

	if len(g.npcList) == 0 {
		g.spawnNPC()
	} else {
		g.npcList[0].update()
	}

	return nil
}

func (g *Game) drawTerrain(screen *ebiten.Image) {

	for y := 0; y*32 < screen.Bounds().Dy(); y += 1 {

		for x := 0; x*32 < screen.Bounds().Dx(); x += 1 {
			op := ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(2), float64(2))

			op.GeoM.Translate(float64(x*32), float64(y*32))
			screen.DrawImage(g.terrainSprites["grass"], &op)
		}
	}

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(2), float64(2))
	op.GeoM.Translate(float64(100), float64(200))
	screen.DrawImage(g.terrainSprites["flowers"], &op)

}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawTerrain(screen)

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(2), float64(2))
	op.GeoM.Translate(float64(g.hero.x), float64(g.hero.y))
	animation := g.hero.activeAnimation

	if g.frameCount%5 == 0 {
		g.hero.nextSprite()
	}

	i := g.hero.sprites[animation][g.hero.currentSprite]
	screen.DrawImage(i, &op)

	for _, npc := range g.npcList {
		if g.frameCount%5 == 0 {
			npc.nextSprite()
		}
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(2), float64(2))
		op.GeoM.Translate(float64(npc.x), float64(npc.y))
		npmSprite := npc.sprites[npc.activeAnimation][npc.currentSprite]
		screen.DrawImage(npmSprite, &op)

	}

	g.frameCount += 1
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) spawnNPC() {
	n := npc{activeAnimation: "down", currentSprite: 0}
	n.init(g)

	g.npcList = append(g.npcList, &n)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Render an image")

	spriteMap := make(map[string][]*ebiten.Image)

	spriteMap["down"] = splitSprites(charSprites, 0, 0, 16, 32, 4)
	spriteMap["right"] = splitSprites(charSprites, 0, 32, 16, 32, 4)
	spriteMap["up"] = splitSprites(charSprites, 0, 64, 16, 32, 4)
	spriteMap["left"] = splitSprites(charSprites, 0, 96, 16, 32, 4)
	spriteMap["idle"] = splitSprites(charSprites, 0, 0, 16, 32, 1)
	spriteMap["hit"] = splitSprites(charSprites, 0, 128, 32, 32, 4)

	if err := ebiten.RunGame(&Game{frameCount: 0, terrainSprites: createTerrainSprites(), hero: character{sprites: spriteMap, activeAnimation: "idle"}}); err != nil {
		log.Fatal(err)
	}
}

func splitSprites(img *ebiten.Image, x int, y int, sizeX int, sizeY int, count int) []*ebiten.Image {
	sprites := []*ebiten.Image{}
	fmt.Printf("splitSprites: x=%v, y=%v, size=%v, count=%v\n", x, y, sizeX, count)

	for i := 0; i < count; i++ {
		min := []int{x + (i * sizeX), y}
		max := []int{x + (i+1)*sizeX, y + sizeY}
		fmt.Printf("Min: %v   ***", min)
		fmt.Printf("Max: %v \n", max)
		sprites = append(sprites,
			img.SubImage(
				image.Rectangle{
					Min: image.Point{X: min[0], Y: min[1]},
					Max: image.Point{X: max[0], Y: max[1]}},
			).(*ebiten.Image))
	}

	return sprites
}

func getSprite(img *ebiten.Image, min []int, max []int) *ebiten.Image {
	i := img.SubImage(
		image.Rectangle{Min: image.Point{X: min[0], Y: min[1]}, Max: image.Point{X: max[0], Y: max[1]}},
	)

	return i.(*ebiten.Image)
}

func createTerrainSprites() map[string]*ebiten.Image {
	terrainSprites, _, err := ebitenutil.NewImageFromFile("assets/overworld.png")
	if err != nil {
		log.Fatal(err)
	}

	terrain := make(map[string]*ebiten.Image)
	terrain["grass"] = getSprite(terrainSprites, []int{0, 0}, []int{16, 16})
	terrain["flowers"] = getSprite(terrainSprites, []int{0, 128}, []int{16, 144})

	return terrain
}
