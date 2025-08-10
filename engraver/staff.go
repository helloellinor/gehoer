package engraver

import (
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Staff struct {
	LengthInStaffSpaces float32
	// You can add thickness or other rendering options here if needed
}

func (s *Staff) Draw(originX, originY float32, fontDefaults map[string]float32) {
	lengthPx := units.StaffSpacesToPixels(s.LengthInStaffSpaces)
	thickness := float32(1)
	if v, ok := fontDefaults["staffLineThickness"]; ok {
		thickness = units.StaffSpacesToPixels(float32(v))
	}

	for i := 0; i < 5; i++ {
		y := originY - units.StaffSpacesToPixels(float32(i))
		rl.DrawLineEx(rl.NewVector2(originX, y), rl.NewVector2(originX+lengthPx, y), thickness, rl.Black)
	}

	// Draw measure bar at the end of staff
	s.DrawMeasureBar(originX+lengthPx, originY, fontDefaults)
}

func (s *Staff) DrawMeasureBar(originX, originY float32, fontDefaults map[string]float32) {
	y1 := originY - units.EmsToPixels(1) // 1 em above bottom line
	y2 := originY
	thickness := float32(2.0)
	if v, ok := fontDefaults["barlineThickness"]; ok {
		thickness = units.StaffSpacesToPixels(float32(v))
	}
	x := originX - thickness/2
	rl.DrawLineEx(rl.NewVector2(x, y1), rl.NewVector2(x, y2), thickness, rl.Black)
}
