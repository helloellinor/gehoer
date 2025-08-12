// Package engraver handles musical notation rendering using SMUFL fonts.
//
// Architecture Overview:
// - Engraver: Main orchestrator that converts music.Score to renderer commands
// - Uses CommandBuffer pattern for efficient batched rendering
// - Supports staff lines, clefs, notes, rests, flags, stems, and accidentals
// - Positioning based on staff spaces with automatic conversions to pixels
//
// Usage:
//
//	engraver := NewEngraver(score, musicFont)
//	engraver.GenerateDrawCommands(x, y, commandBuffer)
//	commandBuffer.Execute(renderer)
package engraver

import (
	"gehoer/music"
	"gehoer/musicfont"
	"gehoer/renderer"
	"gehoer/units"
)

// Engraver converts musical scores into drawing commands using SMUFL fonts
type Engraver struct {
	Score     *music.Score
	MusicFont *musicfont.MusicFont
	fontSize  float32
}

// NewEngraver creates a new engraver for the given score and music font
func NewEngraver(score *music.Score, musicFont *musicfont.MusicFont) *Engraver {
	return &Engraver{
		Score:     score,
		MusicFont: musicFont,
		fontSize:  units.FontRenderSizePx,
	}
}

// CreateGlyphCommand creates a glyph command by name, returning nil if glyph not found
func (e *Engraver) CreateGlyphCommand(glyphName string, x, y float32, color renderer.Color) *renderer.GlyphCommand {
	glyph, ok := e.MusicFont.GetGlyph(glyphName)
	if !ok {
		return nil
	}
	cmd := renderer.NewGlyphCommand(e.MusicFont.Font, glyph.Codepoint, renderer.Vector2{X: x, Y: y}, e.fontSize, color)
	return &cmd
}

// MeasureLengthPx calculates the pixel width needed for a measure
func (e *Engraver) MeasureLengthPx(measure *music.Measure) float32 {
	const elementSpacing = 2.0 // staff spaces between elements
	totalWidth := float32(0)
	spacing := units.StaffSpacesToPixels(elementSpacing)

	for _, elem := range measure.Elements {
		glyphName := elem.GlyphName()
		if glyph, ok := e.MusicFont.GetGlyph(glyphName); ok {
			bboxWidth := e.bboxWidthInPixels(glyph.BBox)
			totalWidth += bboxWidth + spacing
		} else {
			// Fallback width for missing glyphs
			totalWidth += units.StaffSpacesToPixels(2) + spacing
		}
	}
	return totalWidth
}

// Helper to calculate pixel width from GlyphBBox
func (e *Engraver) bboxWidthInPixels(bbox musicfont.GlyphBBox) float32 {
	// bbox.NE[0] = northeast x; bbox.SW[0] = southwest x
	widthStaffSpaces := bbox.NE[0] - bbox.SW[0]
	return units.StaffSpacesToPixels(float32(widthStaffSpaces))
}

// GenerateDrawCommands creates draw commands for the entire score
func (e *Engraver) GenerateDrawCommands(originX, originY float32, buffer *renderer.CommandBuffer) {
	x := originX
	y := originY

	for _, measure := range e.Score.Measures {
		// Calculate staff length for this measure
		staffLength := e.MeasureLengthPx(measure)

		// Add space for clef (always render treble clef for now)
		clefWidth := units.StaffSpacesToPixels(3.0)
		staffLength += clefWidth

		// Draw staff lines
		e.GenerateStaffCommands(x, y, staffLength, renderer.Black, buffer)

		// Draw treble clef at beginning of measure (default)
		if cmd := e.CreateGlyphCommand("gClef", x+units.StaffSpacesToPixels(0.5), y-units.StaffSpacesToPixels(2), renderer.Black); cmd != nil {
			buffer.AddCommand(*cmd)
		}

		// Start drawing elements after clef
		elemX := x + clefWidth

		// Draw measure elements with proportional spacing
		positions := measure.ElementPositions(staffLength-clefWidth, 0, units.StaffSpacesToPixels(1))
		for i, elem := range measure.Elements {
			if i < len(positions) {
				elemX = x + clefWidth + positions[i]
			}

			switch el := elem.(type) {
			case *music.Note:
				e.GenerateNoteCommands(el, elemX, y, renderer.Black, buffer)
			default:
				if cmd := e.CreateGlyphCommand(el.GlyphName(), elemX, y, renderer.Black); cmd != nil {
					buffer.AddCommand(*cmd)
				}
			}
		}

		// Move to next measure with smaller spacing
		x += staffLength + units.StaffSpacesToPixels(2) // reduced from 5 to 2
	}
}
