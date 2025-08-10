package engraver

import (
	"gehoer/music"
	"gehoer/units"

	"gehoer/smufl"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Engraver holds the score, font, and notes drawer.
type Engraver struct {
	Score *music.Score
	Font  rl.Font
	Meta  *smufl.Metadata // Add this field to hold SMuFL metadata
	Notes Notes
	Staff Staff
}

func NewEngraver(score *music.Score, font rl.Font, meta *smufl.Metadata) *Engraver {
	// Calculate staff length based on score content
	staffLength := calculateStaffLength(score)
	
	return &Engraver{
		Score: score,
		Font:  font,
		Meta:  meta,
		Notes: Notes{Score: score, Font: font},
		Staff: Staff{LengthInStaffSpaces: staffLength},
	}
}

func (e *Engraver) Draw(highlightIndex int) {
	staffOriginX := float32(50)
	staffOriginY := float32(200)

	// Draw the staff lines (5 horizontal lines)
	e.Staff.Draw(staffOriginX, staffOriginY, map[string]float32{
		"staffLineThickness": 0.1,
		"barlineThickness":   0.2,
	})

	yTop := staffOriginY - units.StaffSpacesToPixels(4) // 4 staff spaces above bottom line
	yBottom := staffOriginY
	thickness := float32(1.5)

	// Draw vertical barlines at the start of each measure
	for _, measure := range e.Score.Measures {
		measureStartX := staffOriginX + float32((measure.Number-1)*400) // Adjust spacing

		rl.DrawLineEx(rl.NewVector2(measureStartX, yTop), rl.NewVector2(measureStartX, yBottom), thickness, rl.Black)
	}

	// Draw final barline at end of last measure
	lastMeasureX := staffOriginX + float32(len(e.Score.Measures)*400)
	rl.DrawLineEx(rl.NewVector2(lastMeasureX, yTop), rl.NewVector2(lastMeasureX, yBottom), thickness, rl.Black)

	noteCounter := 0

	// Now draw notes in each measure (your existing code)
	for _, measure := range e.Score.Measures {
		measureStartX := staffOriginX + float32((measure.Number-1)*400)

		for i, element := range measure.Elements {
			note, ok := element.(music.Note)
			if !ok {
				continue
			}

			noteX := measureStartX + float32(i)*50
			noteY := staffOriginY - units.StaffSpacesToPixels(float32(note.StaffLine))

			glyphName := note.GetSMUFLName()

			glyphRune, err := e.Meta.GetGlyphRune(glyphName)
			if err != nil {
				continue
			}

			DrawGlyph(e.Font, glyphRune, noteX, noteY, 0, rl.Black)

			if noteCounter == highlightIndex {
				rl.DrawCircleLines(int32(noteX), int32(noteY), 12, rl.Red)
			}

			noteCounter++
		}
	}
}

// Update updates any animation or playback state (currently empty).
func (e *Engraver) Update() {
	// TODO: Add update logic later
}

