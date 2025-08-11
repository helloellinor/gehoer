package engraver

import (
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Staff struct {
	LengthInStaffSpaces float32
	// You can add thickness or other rendering options here if needed
}

func (e *Engraver) DrawStaff(x, y, lengthPx float32, color rl.Color) {
	thickness := float32(e.MusicFont.EngravingDefaults.StaffLineThickness)
	thickness = units.StaffSpacesToPixels(thickness)

	for i := 0; i < 5; i++ {
		// Draw lines stacked vertically with proper spacing
		lineY := y - units.StaffSpacesToPixels(float32(i))
		rl.DrawLineEx(rl.NewVector2(x, lineY), rl.NewVector2(x+lengthPx, lineY), thickness, color)
	}

	// Draw measure barline at the end
	barlineThickness := float32(2.0)

	barlineX := x + lengthPx - barlineThickness/2
	barlineTopY := y - units.StaffSpacesToPixels(4) // top of staff
	barlineBottomY := y

	rl.DrawLineEx(rl.NewVector2(barlineX, barlineTopY), rl.NewVector2(barlineX, barlineBottomY), barlineThickness, color)
}

func (s *Staff) DrawBarline(originX, originY float32, fontDefaults map[string]float32) {
	y1 := originY - units.EmsToPixels(1) // 1 em above bottom line
	y2 := originY
	thickness := float32(2.0)
	if v, ok := fontDefaults["barlineThickness"]; ok {
		thickness = units.StaffSpacesToPixels(float32(v))
	}
	x := originX - thickness/2
	rl.DrawLineEx(rl.NewVector2(x, y1), rl.NewVector2(x, y2), thickness, rl.Black)
}
