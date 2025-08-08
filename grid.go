package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawGrid draws a grid with origin at (originX, originY).
// Coordinates are in Raylib's native system: Y=0 at top, increasing downward.
func DrawGrid(originX, originY float32) {
	// Vertical lines right of origin
	for x := originX; x <= screenWidth; x += staffSpacePx {
		color := rl.LightGray
		if int(x) == int(originX) {
			color = rl.DarkGray
		}
		rl.DrawLine(int32(x), 0, int32(x), screenHeight, color)
		rl.DrawText(fmt.Sprintf("%.0f", x-originX), int32(x)+2, int32(originY)+2, 10, rl.Gray)
	}

	// Vertical lines left of origin
	for x := originX - staffSpacePx; x >= 0; x -= staffSpacePx {
		color := rl.LightGray
		if int(x) == int(originX) {
			color = rl.DarkGray
		}
		rl.DrawLine(int32(x), 0, int32(x), screenHeight, color)
		rl.DrawText(fmt.Sprintf("%.0f", x-originX), int32(x)+2, int32(originY)+2, 10, rl.Gray)
	}

	// Horizontal lines below origin (Y grows downward)
	for y := originY; y <= screenHeight; y += staffSpacePx {
		color := rl.LightGray
		if int(y) == int(originY) {
			color = rl.DarkGray
		}
		rl.DrawLine(0, int32(y), screenWidth, int32(y), color)
		rl.DrawText(fmt.Sprintf("%.0f", y-originY), int32(originX)+2, int32(y)+2, 10, rl.Gray)
	}

	// Horizontal lines above origin (if any)
	for y := originY - staffSpacePx; y >= 0; y -= staffSpacePx {
		color := rl.LightGray
		if int(y) == int(originY) {
			color = rl.DarkGray
		}
		rl.DrawLine(0, int32(y), screenWidth, int32(y), color)
		rl.DrawText(fmt.Sprintf("%.0f", y-originY), int32(originX)+2, int32(y)+2, 10, rl.Gray)
	}
}
