package engraver

import (
	"gehoer/musicfont"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// DebugBBox enables or disables drawing bounding boxes around glyphs globally.
//var DebugBBox = false

// DrawGlyph draws a single SMuFL glyph rune at a baseline origin point
// applying a vertical offset in staff spaces.
// TODO: investigate why offset is in FontRenderSizePx and not in EmPx
// This is useful for drawing glyphs in the context of a staff.
func DrawGlyph(font rl.Font, r rune, originX, originY, verticalOffsetStaffSpaces float32, color rl.Color) {
	verticalOffsetPx := units.StaffSpacesToPixels(verticalOffsetStaffSpaces)

	drawX := originX
	drawY := originY - verticalOffsetPx - (units.FontRenderSizePx / 2)

	rl.DrawTextEx(font, string(r), rl.NewVector2(drawX, drawY), units.FontRenderSizePx, 0.0, color)
	//if DebugBBox {
	//	DrawBBox(metadata.GlyphBBox{}bbox, originX, originY, verticalOffsetStaffSpaces)
	//}
}

// DrawBBox draws a red rectangle around the bounding box of a glyph, useful for debugging.
// This function is likely not correct in its current form as the origin calculation is not clear.
func DrawBBox(bbox musicfont.GlyphBBox, originX, originY, verticalOffsetStaffSpaces float32) {
	widthPx := units.StaffSpacesToPixels(float32(bbox.NE[0]) - float32(bbox.SW[0]))
	heightPx := units.StaffSpacesToPixels(float32(bbox.NE[1]) - float32(bbox.SW[1]))
	verticalOffsetPx := units.StaffSpacesToPixels(verticalOffsetStaffSpaces)

	rectX := originX + units.StaffSpacesToPixels(float32(bbox.SW[0]))
	rectY := originY - units.StaffSpacesToPixels(float32(bbox.NE[1])) - verticalOffsetPx

	rl.DrawRectangleLinesEx(rl.NewRectangle(rectX, rectY, widthPx, heightPx), 1.0, rl.Red)
}
