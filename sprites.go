package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func splitSprites(img *ebiten.Image, x int, y int, sizeX int, sizeY int, count int) []*ebiten.Image {
	sprites := []*ebiten.Image{}

	for i := 0; i < count; i++ {
		min := []int{x + (i * sizeX), y}
		max := []int{x + (i+1)*sizeX, y + sizeY}
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
