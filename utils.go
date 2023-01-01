package main

import (
	"sort"
)

const OFFSET = 20
const DISTANCE = 40

func doesTouch(direction string, source coord, target coord) bool {
	xDiff := source.x - target.x
	yDiff := source.y - target.y
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

func sortNPCs(list []*npc) []*npc {
	sort.Slice(list, func(i, j int) bool {
		return list[i].coord.y+list[i].spriteOffset.y < list[j].coord.y+list[j].spriteOffset.y
	})

	return list
}
