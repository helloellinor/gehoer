//go:build example
// +build example

package main

import (
	"fmt"
	"gehoer/engraver"
	"gehoer/music"
	"gehoer/musicfont"
	"gehoer/renderer"
)

// Example demonstrating the separation of concerns:
// This example shows how the engraver can generate layout commands
// without requiring a graphics context, making it testable.
// Run with: go run -tags example example_separation.go

func main() {
	// Create a simple score
	score := music.NewScore("Test", "Test Composer", "C", "major", 4, 4, 120)
	measure := score.AddMeasure(nil)

	// Add a simple note
	note := &music.Note{
		Pitch:      60, // Middle C
		Duration:   music.QuarterNote,
		StaffLine:  4, // Middle of staff
		Accidental: "",
	}
	measure.AddNote(note)

	// Mock font for demonstration (in real usage, load from file)
	mockFont := &musicfont.MusicFont{
		EngravingDefaults: musicfont.EngravingDefaults{
			StaffLineThickness: 0.1,
			StemThickness:      0.12,
		},
	}

	// Create engraver (layout calculator)
	eng := engraver.NewEngraver(score, mockFont)

	// Create command buffer
	buffer := renderer.NewCommandBuffer()

	// Generate layout commands without any graphics context
	eng.GenerateDrawCommands(100, 200, buffer)

	// At this point, we have all the drawing commands calculated
	// but nothing has been drawn yet. This demonstrates complete
	// separation of layout from rendering.

	fmt.Println("Generated drawing commands successfully!")
	fmt.Println("Layout calculation complete - no graphics context required!")
	fmt.Println("This demonstrates that layout logic is now completely separated from rendering.")

	// In a real application, you would now execute the commands:
	// renderer := renderer.NewRaylibRenderer()
	// buffer.Execute(renderer)
}
