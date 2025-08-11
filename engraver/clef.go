package engraver

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawClef draws a clef glyph at position with color
func (e *Engraver) DrawClef(clefName string, x, y float32, color rl.Color) {
	e.DrawGlyph(clefName, x, y, color)
}
