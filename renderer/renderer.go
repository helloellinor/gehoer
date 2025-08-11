package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Vector2 represents a 2D point to avoid dependency on Raylib in interfaces
type Vector2 struct {
	X, Y float32
}

// Color represents a color to avoid dependency on Raylib in interfaces
type Color struct {
	R, G, B, A uint8
}

// Predefined colors
var (
	Black     = Color{0, 0, 0, 255}
	White     = Color{255, 255, 255, 255}
	Red       = Color{255, 0, 0, 255}
	DarkGray  = Color{80, 80, 80, 255}
	LightGray = Color{220, 220, 220, 255}
)

// DrawCommand represents a drawing operation
type DrawCommand interface {
	Execute(renderer Renderer)
}

// Renderer defines the interface for rendering operations
type Renderer interface {
	DrawLine(start, end Vector2, thickness float32, color Color)
	DrawText(text string, position Vector2, fontSize int32, color Color)
	DrawGlyph(font rl.Font, glyph rune, position Vector2, fontSize float32, color Color)
	DrawRectangleLines(x, y, width, height, lineThickness float32, color Color)
}

// LineCommand represents a line drawing operation
type LineCommand struct {
	Start     Vector2
	End       Vector2
	Thickness float32
	Color     Color
}

func (cmd LineCommand) Execute(renderer Renderer) {
	renderer.DrawLine(cmd.Start, cmd.End, cmd.Thickness, cmd.Color)
}

// TextCommand represents a text drawing operation
type TextCommand struct {
	Text     string
	Position Vector2
	FontSize int32
	Color    Color
}

func (cmd TextCommand) Execute(renderer Renderer) {
	renderer.DrawText(cmd.Text, cmd.Position, cmd.FontSize, cmd.Color)
}

// GlyphCommand represents a glyph drawing operation
type GlyphCommand struct {
	Font     rl.Font
	Glyph    rune
	Position Vector2
	FontSize float32
	Color    Color
}

func (cmd GlyphCommand) Execute(renderer Renderer) {
	renderer.DrawGlyph(cmd.Font, cmd.Glyph, cmd.Position, cmd.FontSize, cmd.Color)
}

// RectangleLinesCommand represents a rectangle outline drawing operation
type RectangleLinesCommand struct {
	X, Y, Width, Height float32
	LineThickness       float32
	Color               Color
}

func (cmd RectangleLinesCommand) Execute(renderer Renderer) {
	renderer.DrawRectangleLines(cmd.X, cmd.Y, cmd.Width, cmd.Height, cmd.LineThickness, cmd.Color)
}

// CommandBuffer accumulates drawing commands
type CommandBuffer struct {
	commands []DrawCommand
}

func NewCommandBuffer() *CommandBuffer {
	return &CommandBuffer{
		commands: make([]DrawCommand, 0),
	}
}

func (cb *CommandBuffer) AddCommand(cmd DrawCommand) {
	cb.commands = append(cb.commands, cmd)
}

func (cb *CommandBuffer) Clear() {
	cb.commands = cb.commands[:0]
}

func (cb *CommandBuffer) Execute(renderer Renderer) {
	for _, cmd := range cb.commands {
		cmd.Execute(renderer)
	}
}

// Helper functions to create commands
func NewLineCommand(start, end Vector2, thickness float32, color Color) LineCommand {
	return LineCommand{Start: start, End: end, Thickness: thickness, Color: color}
}

func NewTextCommand(text string, position Vector2, fontSize int32, color Color) TextCommand {
	return TextCommand{Text: text, Position: position, FontSize: fontSize, Color: color}
}

func NewGlyphCommand(font rl.Font, glyph rune, position Vector2, fontSize float32, color Color) GlyphCommand {
	return GlyphCommand{Font: font, Glyph: glyph, Position: position, FontSize: fontSize, Color: color}
}

func NewRectangleLinesCommand(x, y, width, height, lineThickness float32, color Color) RectangleLinesCommand {
	return RectangleLinesCommand{X: x, Y: y, Width: width, Height: height, LineThickness: lineThickness, Color: color}
}
