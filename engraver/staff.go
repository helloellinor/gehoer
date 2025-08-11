package engraver

import (
	"gehoer/renderer"
	"gehoer/units"
)

type Staff struct {
	LengthInStaffSpaces float32
	// You can add thickness or other rendering options here if needed
}

func (e *Engraver) GenerateStaffCommands(x, y, lengthPx float32, color renderer.Color, buffer *renderer.CommandBuffer) {
	thickness := float32(e.MusicFont.EngravingDefaults.StaffLineThickness)
	thickness = units.StaffSpacesToPixels(thickness)

	for i := 0; i < 5; i++ {
		// Draw lines stacked vertically with proper spacing
		lineY := y - units.StaffSpacesToPixels(float32(i))
		start := renderer.Vector2{X: x, Y: lineY}
		end := renderer.Vector2{X: x + lengthPx, Y: lineY}
		buffer.AddCommand(renderer.NewLineCommand(start, end, thickness, color))
	}

	// Draw measure barline at the end
	barlineThickness := float32(2.0)
	barlineX := x + lengthPx - barlineThickness/2
	barlineTopY := y - units.StaffSpacesToPixels(4) // top of staff
	barlineBottomY := y

	start := renderer.Vector2{X: barlineX, Y: barlineTopY}
	end := renderer.Vector2{X: barlineX, Y: barlineBottomY}
	buffer.AddCommand(renderer.NewLineCommand(start, end, barlineThickness, color))
}

func (s *Staff) GenerateBarlineCommands(originX, originY float32, fontDefaults map[string]float32, buffer *renderer.CommandBuffer) {
	y1 := originY - units.EmsToPixels(1) // 1 em above bottom line
	y2 := originY
	thickness := float32(2.0)
	if v, ok := fontDefaults["barlineThickness"]; ok {
		thickness = units.StaffSpacesToPixels(float32(v))
	}
	x := originX - thickness/2

	start := renderer.Vector2{X: x, Y: y1}
	end := renderer.Vector2{X: x, Y: y2}
	buffer.AddCommand(renderer.NewLineCommand(start, end, thickness, renderer.Black))
}
