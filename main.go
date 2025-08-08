package main

import (
	"log"

	"gehoer/fontloader"
	"gehoer/smufl"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawGlyphWithBBox(
	font rl.Font,
	r rune,
	bbox smufl.GlyphBBox,
	originX, originY float32,
	staffSpacePx float32,
) {
	// 1. Calculate box dimensions in pixels (This was always correct)
	width := (float32(bbox.BBoxNE[0]) - float32(bbox.BBoxSW[0])) * (staffSpacePx / 4)
	height := (float32(bbox.BBoxNE[1]) - float32(bbox.BBoxSW[1])) * (staffSpacePx / 4)

	rl.DrawRectangleLinesEx(rl.NewRectangle(originX, (originY-(height/2)), width, height), 1, rl.Red)

	rl.DrawTextEx(font, string(r), rl.NewVector2(originX, (originY-(fontSize/2))), float32(font.BaseSize), 0, rl.Black)
}

func main() {
	originX := float32(screenWidth / 2)
	originY := float32(screenHeight / 2)
	wholeNoteRune := '\uE0A2'
	wholeNoteName := "noteWhole"

	rl.InitWindow(screenWidth, screenHeight, "Gehør")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	font, err := fontloader.LoadSmuflFontSafe("assets/fonts/Leland/Leland.otf", fontSize, []rune{wholeNoteRune})
	if err != nil {
		log.Fatalf("Font load error: %v", err)
	}
	defer rl.UnloadFont(font)

	smuflData, err := smufl.LoadSmuflMetadata("assets/fonts/Leland/leland_metadata.json")
	if err != nil {
		log.Fatalf("SMuFL metadata load error: %v", err)
	}

	bbox, ok := smuflData.GlyphBBoxes[wholeNoteName]
	if !ok {
		log.Fatalf("Bounding box for glyph %q not found", wholeNoteName)
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw grid (in grid.go)
		DrawGrid(originX, originY)

		// Draw glyph + bounding box using native coords, no Y flip
		drawGlyphWithBBox(font, wholeNoteRune, bbox, originX, originY, staffSpacePx)

		// Draw baseline and origin marker
		rl.DrawLine(int32(originX-50), int32(originY), int32(originX+50), int32(originY), rl.Green)
		rl.DrawCircleV(rl.NewVector2(originX, originY), 5, rl.Blue)
		rl.DrawCircle(0, 0, 5, rl.Red)
		rl.DrawText("Top-left (0,0)", 10, 10, 10, rl.Red)

		rl.EndDrawing()
	}
}
