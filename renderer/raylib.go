package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// RaylibRenderer implements the Renderer interface using Raylib
type RaylibRenderer struct{}

func NewRaylibRenderer() *RaylibRenderer {
	return &RaylibRenderer{}
}

// Helper function to convert our Color to Raylib Color
func (r *RaylibRenderer) toRaylibColor(color Color) rl.Color {
	return rl.Color{R: color.R, G: color.G, B: color.B, A: color.A}
}

// Helper function to convert our Vector2 to Raylib Vector2
func (r *RaylibRenderer) toRaylibVector2(v Vector2) rl.Vector2 {
	return rl.Vector2{X: v.X, Y: v.Y}
}

func (r *RaylibRenderer) DrawLine(start, end Vector2, thickness float32, color Color) {
	rl.DrawLineEx(r.toRaylibVector2(start), r.toRaylibVector2(end), thickness, r.toRaylibColor(color))
}

func (r *RaylibRenderer) DrawText(text string, position Vector2, fontSize int32, color Color) {
	rl.DrawText(text, int32(position.X), int32(position.Y), fontSize, r.toRaylibColor(color))
}

func (r *RaylibRenderer) DrawGlyph(font rl.Font, glyph rune, position Vector2, fontSize float32, color Color) {
	rl.DrawTextEx(font, string(glyph), r.toRaylibVector2(position), fontSize, 0, r.toRaylibColor(color))
}

func (r *RaylibRenderer) DrawRectangleLines(x, y, width, height, lineThickness float32, color Color) {
	rl.DrawRectangleLinesEx(rl.NewRectangle(x, y, width, height), lineThickness, r.toRaylibColor(color))
}
