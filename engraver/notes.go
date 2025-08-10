package engraver

import (
	"gehoer/music"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Notes handles drawing noteheads and rests.
type Notes struct {
	Score *music.Score
	Font  rl.Font
}

// Draw draws all notes/rests in the score, highlighting one note by index.
func (n *Notes) Draw(highlightIndex int) {
	if n.Score == nil || len(n.Score.Measures) == 0 {
		return
	}

	startX := float32(50)
	startY := float32(200)

	noteCounter := 0

	for _, measure := range n.Score.Measures {
		positions := measure.ElementPositions(units.EmPx*8, units.EmPx*0.5, units.EmPx*0.5)
		for i, element := range measure.Elements {
			x := startX + positions[i]
			y := startY

			color := rl.Black
			if noteCounter == highlightIndex {
				color = rl.Red
			}

			switch elem := element.(type) {
			case music.Note:
				n.drawNotehead(elem, x, y, color)
			case music.Rest:
				n.drawRest(elem, x, y, color)
			}

			noteCounter++
		}
		startY += units.StaffSpacePx * 7
	}
}

func (n *Notes) drawNotehead(note music.Note, x, y float32, color rl.Color) {
	var runeToDraw rune
	switch note.Duration {
	case music.WholeNote:
		runeToDraw = '\uE0A2' // whole note
	case music.HalfNote:
		runeToDraw = '\uE0A3' // half note
	default:
		runeToDraw = '\uE0A4' // quarter note default
	}

	offset := float32(note.StaffLine)

	DrawGlyph(n.Font, runeToDraw, x, y, offset, color)
}

func (n *Notes) drawRest(rest music.Rest, x, y float32, color rl.Color) {
	rl.DrawText("Rest", int32(x), int32(y), 12, color)
}
