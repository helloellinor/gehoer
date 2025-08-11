package engraver

import (
	"gehoer/musicfont"
	"gehoer/renderer"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// GlyphDrawInfo contains positioning information for drawing a glyph
type GlyphDrawInfo struct {
	Font                      rl.Font
	Glyph                     rune
	OriginX, OriginY          float32
	VerticalOffsetStaffSpaces float32
	FontSize                  float32
	Color                     renderer.Color
}

// CalculateGlyphPosition calculates the final draw position for a glyph
// applying a vertical offset in staff spaces.
func CalculateGlyphPosition(originX, originY, verticalOffsetStaffSpaces float32) renderer.Vector2 {
	verticalOffsetPx := units.StaffSpacesToPixels(verticalOffsetStaffSpaces)
	drawX := originX
	drawY := originY - verticalOffsetPx - (units.FontRenderSizePx / 2)
	return renderer.Vector2{X: drawX, Y: drawY}
}

// CreateGlyphCommand creates a glyph draw command
func CreateGlyphCommand(font rl.Font, glyph rune, originX, originY, verticalOffsetStaffSpaces float32, color renderer.Color) renderer.GlyphCommand {
	position := CalculateGlyphPosition(originX, originY, verticalOffsetStaffSpaces)
	return renderer.NewGlyphCommand(font, glyph, position, units.FontRenderSizePx, color)
}

// CalculateBBoxPosition calculates the bounding box rectangle for debugging
func CalculateBBoxPosition(bbox musicfont.GlyphBBox, originX, originY, verticalOffsetStaffSpaces float32) (x, y, width, height float32) {
	widthPx := units.StaffSpacesToPixels(float32(bbox.NE[0]) - float32(bbox.SW[0]))
	heightPx := units.StaffSpacesToPixels(float32(bbox.NE[1]) - float32(bbox.SW[1]))
	verticalOffsetPx := units.StaffSpacesToPixels(verticalOffsetStaffSpaces)

	rectX := originX + units.StaffSpacesToPixels(float32(bbox.SW[0]))
	rectY := originY - units.StaffSpacesToPixels(float32(bbox.NE[1])) - verticalOffsetPx

	return rectX, rectY, widthPx, heightPx
}

// CreateBBoxCommand creates a bounding box debug draw command
func CreateBBoxCommand(bbox musicfont.GlyphBBox, originX, originY, verticalOffsetStaffSpaces float32) renderer.RectangleLinesCommand {
	x, y, width, height := CalculateBBoxPosition(bbox, originX, originY, verticalOffsetStaffSpaces)
	return renderer.NewRectangleLinesCommand(x, y, width, height, 1.0, renderer.Red)
}
