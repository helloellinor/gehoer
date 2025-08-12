// Package engraver - glyph utilities for positioning and rendering music glyphs
package engraver

import (
	"gehoer/renderer"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// CreateGlyphCommand creates a glyph draw command with optional vertical offset
func CreateGlyphCommand(font rl.Font, glyph rune, originX, originY, verticalOffsetStaffSpaces float32, color renderer.Color) renderer.GlyphCommand {
	verticalOffsetPx := units.StaffSpacesToPixels(verticalOffsetStaffSpaces)
	position := renderer.Vector2{
		X: originX,
		Y: originY - verticalOffsetPx,
	}
	return renderer.NewGlyphCommand(font, glyph, position, units.FontRenderSizePx, color)
}
