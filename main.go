package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var charSprites *ebiten.Image

type coord struct {
	x int
	y int
}

const MOVE_DELTA = 1

var hitAnimationMap = map[string]string{
	"right": "hitRight",
	"left":  "hitLeft",
	"up":    "hitUp",
	"down":  "hitDown",
}

func init() {

	var err error
	charSprites, _, err = ebitenutil.NewImageFromFile("assets/character.png")

	if err != nil {
		log.Fatal(err)
	}

}

type Game struct {
	hero           character
	frameCount     int
	terrainSprites map[string]*ebiten.Image

	npcList []*npc
	objects []*object

	dialogs       map[string][]dialog
	dialogChain   []dialog
	currentDialog int

	debugPoints []coord
	debug       bool

	initialized bool
}

func (g *Game) init() {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(2), float64(2))
	g.createObject(g.terrainSprites["flowers"], coord{100, 100}, coord{-16, -16}, dialogHandlerFactory("flower_look"))
	g.initialized = true
}

func (g *Game) Update() error {
	if !g.initialized {
		g.init()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.dialogChain != nil {
			if g.currentDialog+2 > len(g.dialogChain) {
				g.dialogChain = nil
				g.currentDialog = 0
				return nil
			} else {
				g.currentDialog += 1
				return nil
			}
		}
		g.hero.hitting = true
		fmt.Println("action is triggered")
		g.checkCollisions()
		g.hero.actionCounter = 30
		g.hero.activeAnimation = hitAnimationMap[g.hero.direction]
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.hero.walking = true
		g.hero.activeAnimation = "left"
		g.hero.direction = "left"
		g.hero.coord.x -= MOVE_DELTA
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.hero.walking = true
		g.hero.activeAnimation = "right"
		g.hero.direction = "right"
		g.hero.coord.x += MOVE_DELTA
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.hero.walking = true
		g.hero.activeAnimation = "down"
		g.hero.direction = "down"
		g.hero.coord.y += MOVE_DELTA
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.hero.walking = true
		g.hero.activeAnimation = "up"
		g.hero.direction = "up"
		g.hero.coord.y -= MOVE_DELTA
	} else if g.hero.actionCounter > 0 {
		g.hero.actionCounter -= 1
	} else {
		g.hero.currentSprite = 0
		g.hero.walking = false
		g.hero.hitting = false
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

}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawTerrain(screen)

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(2), float64(2))
	op.GeoM.Translate(float64(g.hero.coord.x+g.hero.spriteOffset.x), float64(g.hero.coord.y+g.hero.spriteOffset.y))
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
		op.GeoM.Translate(float64(npc.coord.x), float64(npc.coord.y))
		npmSprite := npc.sprites[npc.activeAnimation][npc.currentSprite]
		screen.DrawImage(npmSprite, &op)
	}

	for _, obj := range g.objects {
		x := obj.coord.x + obj.spriteOffset.x
		y := obj.coord.y + obj.spriteOffset.y
		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(2), float64(2))
		op.GeoM.Translate(float64(x), float64(y))

		screen.DrawImage(obj.sprite, &op)
	}

	g.frameCount += 1
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	mplusNormalFont, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if g.dialogChain != nil {
		d := (g.dialogChain)[g.currentDialog]
		ebitenutil.DrawRect(screen, 0, float64(screen.Bounds().Dy()/5*4), float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy()/5*2), color.White)
		text.Draw(screen, d.Text, mplusNormalFont, 50, screen.Bounds().Dy()/5*4, color.Black)
	}

	if g.debug {
		for _, dp := range g.objects {
			ebitenutil.DrawCircle(screen, float64(dp.coord.x), float64(dp.coord.y), 5, color.Black)
		}
		ebitenutil.DrawCircle(screen, float64(g.hero.coord.x), float64(g.hero.coord.y), 5, color.Black)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func (g *Game) spawnNPC() {
	n := npc{activeAnimation: "down", currentSprite: 0}
	n.init(g)

	g.npcList = append(g.npcList, &n)
}

func (g *Game) checkCollisions() *object {
	for _, object := range g.objects {
		touch := doesTouch(g.hero.direction, g.hero.coord, object.coord)
		if touch {
			if object.handler != nil {
				object.handler(g)
			}
			return object
		}
	}
	return nil
}

func (g *Game) createObject(sprite *ebiten.Image, coord coord, spriteOffset coord, handler func(*Game)) {
	obj := object{coord: coord, spriteOffset: spriteOffset, sprite: sprite, handler: handler}
	g.objects = append(g.objects, &obj)
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
	spriteMap["hitDown"] = splitSprites(charSprites, 7, 128, 32, 32, 4)
	spriteMap["hitUp"] = splitSprites(charSprites, 7, 160, 32, 32, 4)
	spriteMap["hitRight"] = splitSprites(charSprites, 7, 192, 32, 32, 4)
	spriteMap["hitLeft"] = splitSprites(charSprites, 7, 224, 32, 32, 4)

	data, err := os.ReadFile("dialogs/start_scene.json")
	if err != nil {
		panic(err)
	}
	dialogs := loadDialogs(data)

	if err := ebiten.RunGame(
		&Game{
			initialized:    false,
			frameCount:     0,
			terrainSprites: createTerrainSprites(),
			dialogChain:    dialogs["opening"],
			hero: character{
				sprites: spriteMap, activeAnimation: "idle", direction: "down", spriteOffset: coord{x: -16, y: -32},
			},
			dialogs: dialogs,
			debug:   true,
		}); err != nil {
		log.Fatal(err)
	}
}
