// Package engraver - staff rendering utilities
package engraver

import (
	"gehoer/renderer"
	"gehoer/units"
)

// GenerateStaffCommands draws the five staff lines and a barline at the end
func (e *Engraver) GenerateStaffCommands(x, y, lengthPx float32, color renderer.Color, buffer *renderer.CommandBuffer) {
	thickness := units.StaffSpacesToPixels(float32(e.MusicFont.EngravingDefaults.StaffLineThickness))

	// Draw the five staff lines
	for i := 0; i < 5; i++ {
		lineY := y - units.StaffSpacesToPixels(float32(i))
		start := renderer.Vector2{X: x, Y: lineY}
		end := renderer.Vector2{X: x + lengthPx, Y: lineY}
		buffer.AddCommand(renderer.NewLineCommand(start, end, thickness, color))
	}

	// Draw barline at the end
	barlineThickness := float32(2.0)
	barlineX := x + lengthPx - barlineThickness/2
	start := renderer.Vector2{X: barlineX, Y: y - units.StaffSpacesToPixels(4)}
	end := renderer.Vector2{X: barlineX, Y: y}
	buffer.AddCommand(renderer.NewLineCommand(start, end, barlineThickness, color))
}
