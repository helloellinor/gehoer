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
	cmd := CreateGlyphCommand(e.MusicFont.Font, glyph.Codepoint, x, y, 0, color)
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

// GenerateDrawCommands creates draw commands for the entire score with line breaking
func (e *Engraver) GenerateDrawCommands(originX, originY float32, buffer *renderer.CommandBuffer) {
	const maxLineWidth = 1000.0 // Maximum width before line break (pixels)
	const lineHeight = 150.0    // Height between staff lines (pixels)

	x := originX
	y := originY
	lineWidth := float32(0)

	for _, measure := range e.Score.Measures {
		// Calculate staff length for this measure
		staffLength := e.MeasureLengthPx(measure)

		// Add space for clef and key signature (only for first measure or after line break)
		clefWidth := units.StaffSpacesToPixels(3.0)
		keySignatureWidth := e.getKeySignatureWidth()
		prefixWidth := clefWidth + keySignatureWidth

		// Check if this is the first measure of a line
		isFirstMeasureOfLine := lineWidth == 0
		if isFirstMeasureOfLine {
			staffLength += prefixWidth
		}

		measureSpacing := units.StaffSpacesToPixels(2)
		totalMeasureWidth := staffLength + measureSpacing

		// Check if we need to break to a new line
		if lineWidth > 0 && lineWidth+totalMeasureWidth > maxLineWidth {
			// Start new line
			x = originX
			y += lineHeight
			lineWidth = 0
			isFirstMeasureOfLine = true
			staffLength += prefixWidth // Add prefix for new line
			totalMeasureWidth = staffLength + measureSpacing
		}

		// Draw staff lines
		e.GenerateStaffCommands(x, y, staffLength, renderer.Black, buffer)

		// Draw clef and key signature if first measure of line
		if isFirstMeasureOfLine {
			// Draw treble clef
			if cmd := e.CreateGlyphCommand("gClef", x+units.StaffSpacesToPixels(0.5), y-units.StaffSpacesToPixels(2), renderer.Black); cmd != nil {
				buffer.AddCommand(*cmd)
			}

			// Draw key signature after clef
			e.GenerateKeySignatureCommands(x+clefWidth, y, renderer.Black, buffer)
		}

		// Start drawing elements
		elemStartX := x
		if isFirstMeasureOfLine {
			elemStartX += prefixWidth
		}

		// Draw measure elements with proportional spacing
		availableWidth := staffLength
		if isFirstMeasureOfLine {
			availableWidth -= prefixWidth
		}
		positions := measure.ElementPositions(availableWidth, 0, units.StaffSpacesToPixels(1))
		for i, elem := range measure.Elements {
			elemX := elemStartX
			if i < len(positions) {
				elemX = elemStartX + positions[i]
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

		// Move to next measure position
		x += totalMeasureWidth
		lineWidth += totalMeasureWidth
	}
}

// getKeySignatureWidth calculates the width needed for the key signature
func (e *Engraver) getKeySignatureWidth() float32 {
	sharps, flats := e.getKeySignatureAccidentals()
	if sharps == 0 && flats == 0 {
		return 0 // C major/A minor has no accidentals
	}
	accidentalCount := sharps + flats
	return units.StaffSpacesToPixels(float32(accidentalCount) * 0.75) // 0.75 staff spaces per accidental
}

// GenerateKeySignatureCommands draws the key signature accidentals
func (e *Engraver) GenerateKeySignatureCommands(x, y float32, color renderer.Color, buffer *renderer.CommandBuffer) {
	sharps, flats := e.getKeySignatureAccidentals()

	if sharps > 0 {
		e.drawAccidentals(x, y, sharps, true, color, buffer)
	} else if flats > 0 {
		e.drawAccidentals(x, y, flats, false, color, buffer)
	}
}

// getKeySignatureAccidentals returns the number of sharps or flats for the current key
func (e *Engraver) getKeySignatureAccidentals() (sharps, flats int) {
	tonic := e.Score.KeySignature.Tonic
	mode := e.Score.KeySignature.Mode

	// Major keys (only supporting major keys for now)
	if mode == "dur" || mode == "major" {
		majorKeyMap := map[string][2]int{
			"C": {0, 0}, "G": {1, 0}, "D": {2, 0}, "A": {3, 0}, "E": {4, 0}, "B": {5, 0}, "F#": {6, 0}, "C#": {7, 0},
			"F": {0, 1}, "Bb": {0, 2}, "Eb": {0, 3}, "Ab": {0, 4}, "Db": {0, 5}, "Gb": {0, 6}, "Cb": {0, 7},
		}
		if accidentals, ok := majorKeyMap[tonic]; ok {
			return accidentals[0], accidentals[1]
		}
	}

	return 0, 0 // Default to C major
}

// drawAccidentals draws sharps or flats in key signature order
func (e *Engraver) drawAccidentals(x, y float32, count int, isSharp bool, color renderer.Color, buffer *renderer.CommandBuffer) {
	var glyphName string
	var positions []float32

	if isSharp {
		glyphName = "accidentalSharp"
		positions = []float32{2.5, 3.5, 2, 3, 1.5, 2.5, 1} // F#, C#, G#, D#, A#, E#, B#
	} else {
		glyphName = "accidentalFlat"
		positions = []float32{3, 1.5, 3.5, 2, 4, 2.5, 4.5} // Bb, Eb, Ab, Db, Gb, Cb, Fb
	}

	for i := 0; i < count && i < len(positions); i++ {
		accX := x + units.StaffSpacesToPixels(float32(i)*0.75)
		accY := y - units.StaffSpacesToPixels(positions[i])

		if cmd := e.CreateGlyphCommand(glyphName, accX, accY, color); cmd != nil {
			buffer.AddCommand(*cmd)
		}
	}
}
