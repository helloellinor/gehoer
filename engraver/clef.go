package engraver

import (
	"gehoer/renderer"
	"gehoer/units"
)

// GenerateClefCommands creates a draw command for a clef glyph at position with color
func (e *Engraver) GenerateClefCommands(clefName string, x, y float32, color renderer.Color, buffer *renderer.CommandBuffer) {
	glyph, ok := e.MusicFont.GetGlyph(clefName)
	if !ok {
		return
	}

	// Apply baseline adjustment for proper clef positioning on staff
	// SMUFL coordinates: Y=0 is staff line, Y increases upward (mathematical coords)
	// Screen coordinates: Y increases downward 
	// If glyph center is above SMUFL baseline (positive), move glyph down in screen coords (add pixels)
	// If glyph center is below SMUFL baseline (negative), move glyph up in screen coords (subtract pixels)
	glyphCenterY := (glyph.BBox.SW[1] + glyph.BBox.NE[1]) / 2
	baselineAdjustment := units.StaffSpacesToPixels(float32(glyphCenterY))
	adjustedY := y + baselineAdjustment

	cmd := CreateGlyphCommand(e.MusicFont.Font, glyph.Codepoint, x, adjustedY, 0, color)
	buffer.AddCommand(cmd)
}
