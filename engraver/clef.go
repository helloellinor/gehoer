package engraver

import (
	"gehoer/renderer"
)

// GenerateClefCommands creates a draw command for a clef glyph at position with color
func (e *Engraver) GenerateClefCommands(clefName string, x, y float32, color renderer.Color, buffer *renderer.CommandBuffer) {
	if cmd := e.CreateGlyphCommand(clefName, x, y, color); cmd != nil {
		buffer.AddCommand(*cmd)
	}
}
