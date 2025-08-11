package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Utility functions to convert between Raylib types and renderer types

// FromRaylibColor converts a Raylib color to renderer color
func FromRaylibColor(rlColor rl.Color) Color {
	return Color{R: rlColor.R, G: rlColor.G, B: rlColor.B, A: rlColor.A}
}

// ToRaylibColor converts a renderer color to Raylib color
func ToRaylibColor(color Color) rl.Color {
	return rl.Color{R: color.R, G: color.G, B: color.B, A: color.A}
}

// FromRaylibVector2 converts a Raylib Vector2 to renderer Vector2
func FromRaylibVector2(rlVec rl.Vector2) Vector2 {
	return Vector2{X: rlVec.X, Y: rlVec.Y}
}

// ToRaylibVector2 converts a renderer Vector2 to Raylib Vector2
func ToRaylibVector2(vec Vector2) rl.Vector2 {
	return rl.Vector2{X: vec.X, Y: vec.Y}
}

// Common Raylib colors converted to renderer colors
var (
	RayWhite = FromRaylibColor(rl.RayWhite)
	Gray     = FromRaylibColor(rl.Gray)
)
