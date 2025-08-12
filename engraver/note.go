package engraver

import (
	"gehoer/music"
	"gehoer/renderer"
	"gehoer/units"
)

// GenerateNoteCommands creates all drawing commands for a note:
// - Notehead positioned by staff line
// - Stem (if not whole note) with proper direction and length
// - Flags (for eighth notes and shorter) with correct rotation
// - Accidentals (if present) positioned to the left
// - Ledger lines (if note is outside normal staff range)
func (e *Engraver) GenerateNoteCommands(note *music.Note, x, y float32, color renderer.Color, buffer *renderer.CommandBuffer) {
	noteheadName := note.NoteheadGlyphName()
	glyph, ok := e.MusicFont.GetGlyph(noteheadName)
	if !ok {
		return
	}

	// Vertical offset for the notehead position on staff lines (staffLine = 0 is bottom line)
	verticalOffsetPx := units.StaffSpacesToPixels(float32(note.StaffLine))
	noteheadX := x
	noteheadY := y - verticalOffsetPx

	// Draw notehead
	cmd := CreateGlyphCommand(e.MusicFont.Font, glyph.Codepoint, noteheadX, noteheadY, 0, color)
	buffer.AddCommand(cmd)

	if note.Duration != music.WholeNote {
		stemUp := note.StaffLine < 3 // example rule for stem direction
		stemLength := units.StaffSpacesToPixels(3.5)
		stemThickness := units.StaffSpacesToPixels(float32(e.MusicFont.EngravingDefaults.StemThickness))

		var stemX, stemStartY, stemEndY float32

		if stemUp {
			if a, ok := e.MusicFont.Anchors[noteheadName]["stemUpSE"]; ok {
				stemX = noteheadX + units.StaffSpacesToPixels(float32(a[0]))
				stemStartY = noteheadY - units.StaffSpacesToPixels(float32(a[1]))
			} else {
				// fallback if no anchor
				stemX = noteheadX + units.StaffSpacesToPixels(0.5)
				stemStartY = noteheadY
			}
			stemEndY = stemStartY - stemLength
		} else {
			if a, ok := e.MusicFont.Anchors[noteheadName]["stemDownNW"]; ok {
				stemX = noteheadX + units.StaffSpacesToPixels(float32(a[0]))
				stemStartY = noteheadY - units.StaffSpacesToPixels(float32(a[1]))
			} else {
				stemX = noteheadX - units.StaffSpacesToPixels(0.5)
				stemStartY = noteheadY
			}
			stemEndY = stemStartY + stemLength
		}

		stemDrawX := stemX
		if stemUp {
			stemDrawX -= stemThickness / 2
		} else {
			stemDrawX += stemThickness / 2
		}

		// Draw stem
		start := renderer.Vector2{X: stemDrawX, Y: stemStartY}
		end := renderer.Vector2{X: stemDrawX, Y: stemEndY}
		buffer.AddCommand(renderer.NewLineCommand(start, end, stemThickness, color))

		// Draw flags for eighth notes and shorter
		if note.HasFlag() {
			flagName := e.flagGlyphName(note.Duration, stemUp)
			if flagGlyph, ok := e.MusicFont.GetGlyph(flagName); ok {
				// Simplified flag positioning
				flagX := stemDrawX
				flagY := stemEndY

				// Adjust position based on stem direction
				if stemUp {
					// Flags attach to the top of stems for stem-up notes
					flagX += units.StaffSpacesToPixels(0.1) // slight right offset
				} else {
					// Flags attach to the bottom of stems for stem-down notes
					flagX -= units.StaffSpacesToPixels(0.1) // slight left offset
				}

				flagCmd := CreateGlyphCommand(e.MusicFont.Font, flagGlyph.Codepoint, flagX, flagY, 0, color)
				buffer.AddCommand(flagCmd)
			}
		}
	}

	// Draw accidental if present
	if note.Accidental != "" {
		accidentalX := x - units.StaffSpacesToPixels(1.5)
		accidentalName := accidentalToGlyphName(note.Accidental)
		accGlyph, ok := e.MusicFont.GetGlyph(accidentalName)
		if ok {
			accCmd := CreateGlyphCommand(e.MusicFont.Font, accGlyph.Codepoint, accidentalX, y, 0, color)
			buffer.AddCommand(accCmd)
		}
	}

	// Draw ledger lines if note is outside staff range
	const bottomStaffLine = 0
	const topStaffLine = 8

	thickness := float32(e.MusicFont.EngravingDefaults.StaffLineThickness)
	thickness = units.StaffSpacesToPixels(thickness)

	centerX := x + units.StaffSpacesToPixels(float32((glyph.BBox.SW[0]+glyph.BBox.NE[0])/2))

	if note.StaffLine < bottomStaffLine {
		// ledger lines below staff
		for line := bottomStaffLine - 2; line >= note.StaffLine; line -= 2 {
			yLine := y + units.StaffSpacesToPixels(float32(bottomStaffLine-line)*0.5)
			start := renderer.Vector2{X: centerX - units.StaffSpacesToPixels(0.75), Y: yLine}
			end := renderer.Vector2{X: centerX + units.StaffSpacesToPixels(0.75), Y: yLine}
			buffer.AddCommand(renderer.NewLineCommand(start, end, thickness, color))
		}
	} else if note.StaffLine > topStaffLine {
		// ledger lines above staff
		for line := topStaffLine + 2; line <= note.StaffLine; line += 2 {
			yLine := y - units.StaffSpacesToPixels(float32(line-topStaffLine)*0.5)
			start := renderer.Vector2{X: centerX - units.StaffSpacesToPixels(0.75), Y: yLine}
			end := renderer.Vector2{X: centerX + units.StaffSpacesToPixels(0.75), Y: yLine}
			buffer.AddCommand(renderer.NewLineCommand(start, end, thickness, color))
		}
	}
}

// Helper function to map accidentals to SMuFL glyph names
func accidentalToGlyphName(acc string) string {
	switch acc {
	case "sharp":
		return "accidentalSharp"
	case "flat":
		return "accidentalFlat"
	case "natural":
		return "accidentalNatural"
	default:
		return ""
	}
}

// Example flagGlyphName helper to get flag glyph by duration and stem direction
func (e *Engraver) flagGlyphName(duration music.NoteValue, stemUp bool) string {
	if stemUp {
		switch duration {
		case music.EighthNote:
			return "flag8thUp"
		case music.SixteenthNote:
			return "flag16thUp"
		case music.ThirtySecondNote:
			return "flag32ndUp"
		case music.SixtyFourthNote:
			return "flag64thUp"
		default:
			return ""
		}
	} else {
		switch duration {
		case music.EighthNote:
			return "flag8thDown"
		case music.SixteenthNote:
			return "flag16thDown"
		case music.ThirtySecondNote:
			return "flag32ndDown"
		case music.SixtyFourthNote:
			return "flag64thDown"
		default:
			return ""
		}
	}
}
