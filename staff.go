package main

import (
	"gehoer/metadata"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// global to allow staff helpers to access renderer defaults
var currentRenderer *ScoreRenderer

// DrawStaff draws 5 staff lines starting from the bottom line upward.
// originX, originY = left end of the bottom staff line
// lengthInStaffSpaces = staff line length in SMUFL staff space units
// In SMUFL: staff lines are separated by exactly 1 staff space
func DrawStaff(originX, originY float32, lengthInStaffSpaces float32) {
	lengthPx := StaffSpacesToPixels(lengthInStaffSpaces)
	thickness := float32(1)
	if currentRenderer != nil && currentRenderer.defaults != nil {
		if v, ok := currentRenderer.defaults["staffLineThickness"]; ok {
			thickness = StaffSpacesToPixels(float32(v))
		}
	}
	for i := 0; i < 5; i++ {
		// Each staff line is exactly 1 staff space above the previous one
		y := originY - StaffSpacesToPixels(float32(i))
		rl.DrawLineEx(rl.NewVector2(originX, y), rl.NewVector2(originX+lengthPx, y), thickness, rl.Black)
	}
}

// DrawMeasureBar draws a measure bar with 1em height at the given position
// originX, originY = left end of the bottom staff line
func DrawMeasureBar(originX, originY float32) {
	// Draw measure bar with height of 1em
	y1 := originY - EmsToPixels(1) // 1 em up from bottom line
	y2 := originY                  // Bottom line
	thickness := float32(2.0)
	if currentRenderer != nil && currentRenderer.defaults != nil {
		if v, ok := currentRenderer.defaults["barlineThickness"]; ok {
			thickness = StaffSpacesToPixels(float32(v))
		}
	}
	// Shift inward by half thickness to keep outer edge inside measure bounds
	x := originX - thickness/2
	rl.DrawLineEx(rl.NewVector2(x, y1), rl.NewVector2(x, y2), thickness, rl.Black)
}

// DrawClef draws a clef symbol at the given staff position.
// clefType: "treble" or "bass"
// originX, originY: left end of the bottom staff line
func DrawClef(clefType string, originX, originY float32, font rl.Font, smufl *metadata.SMuFLMetadata, bboxMap map[string]metadata.GlyphBBox) {
	var clefName string
	switch clefType {
	case "bass":
		clefName = "fClef"
	default:
		clefName = "gClef"
	}
	clefRune, err := smufl.GetGlyphRune(clefName)
	if err != nil {
		return // fallback: do nothing if glyph not found
	}

	// Use bboxMap for bounding box lookup, just like notes/rests
	if bbox, hasBBox := bboxMap[clefName]; hasBBox {
		drawGlyphWithBBox(font, clefRune, bbox, originX+StaffSpacesToPixels(0.5), originY-float32(bbox.NE[1]*2), 0)
	} else {
		rl.DrawTextEx(font, string(clefRune), rl.NewVector2(originX+StaffSpacesToPixels(0.5), originY), FontRenderSize, 0.0, rl.Black)
	}
}
