package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var charSprites *ebiten.Image

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

	npcList         []*npc
	callBackTrigger func(*Game)

	showDialog bool
	text       string
}

const MOVE_DELTA = 1

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.callBackTrigger != nil {
			g.callBackTrigger(g)
			return nil
		}

		fmt.Println("action is triggered")
		g.hero.actionCounter = 30
		g.hero.activeAnimation = "hit"
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
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
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	mplusNormalFont, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if g.showDialog {
		ebitenutil.DrawRect(screen, 0, float64(screen.Bounds().Dy())/4*3-50, float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())/4, color.White)
		text.Draw(screen, "The Border Gate", mplusNormalFont, screen.Bounds().Dx()/3, screen.Bounds().Dy()/4*3, color.Black)
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

	callback := func(game *Game) {
		fmt.Println("you are now calling the action callback :)")
		game.callBackTrigger = nil
		game.showDialog = false
	}

	if err := ebiten.RunGame(&Game{showDialog: true, frameCount: 0, terrainSprites: createTerrainSprites(), callBackTrigger: callback, hero: character{sprites: spriteMap, activeAnimation: "idle"}}); err != nil {
		log.Fatal(err)
	}
}
