package main

import "fmt"

const OFFSET = 20
const DISTANCE = 30

func doesTouch(direction string, source coord, target coord) bool {
	fmt.Println("Does Touch?")
	xDiff := source.x - target.x
	yDiff := source.y - target.y
	fmt.Println("X:", xDiff)
	fmt.Println("Y:", yDiff)

	switch direction {
	case "down":
		return xDiff < OFFSET && xDiff > -(OFFSET) &&
			yDiff < 0 && yDiff > -(DISTANCE)
	case "up":
		return xDiff < OFFSET && xDiff > -(OFFSET) &&
			yDiff > 0 && yDiff < DISTANCE
	case "right":
		return xDiff < 0 && xDiff > -(DISTANCE) &&
			yDiff < OFFSET && yDiff > -(OFFSET)
	case "left":
		return xDiff > 0 && xDiff < DISTANCE &&
			yDiff < OFFSET && yDiff > -(OFFSET)
	}
	return false
}
