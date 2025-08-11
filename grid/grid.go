package grid

import (
	"fmt"
	"gehoer/renderer"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Grid struct {
	Spacing       int            // distance in pixels between lines in world space
	HalfWidth     int            // half width of grid in world units
	HalfHeight    int            // half height of grid in world units
	OriginX       float32        // grid origin X coordinate in world space
	OriginY       float32        // grid origin Y coordinate in world space
	LabelColor    renderer.Color // color for grid labels
	LineColor     renderer.Color // color for grid lines
	LabelFontSize int32          // font size for labels
}

func New(spacing, halfWidth, halfHeight int, labelFontSize int32) *Grid {
	return &Grid{
		Spacing:       spacing,
		HalfWidth:     halfWidth,
		HalfHeight:    halfHeight,
		OriginX:       0,
		OriginY:       0,
		LabelColor:    renderer.DarkGray,
		LineColor:     renderer.LightGray,
		LabelFontSize: labelFontSize,
	}
}

// GenerateDrawCommands creates draw commands for the grid
func (g *Grid) GenerateDrawCommands(camera rl.Camera2D, buffer *renderer.CommandBuffer) {
	// Draw central axes
	start := renderer.Vector2{X: g.OriginX, Y: float32(-g.HalfHeight)}
	end := renderer.Vector2{X: g.OriginX, Y: float32(g.HalfHeight)}
	buffer.AddCommand(renderer.NewLineCommand(start, end, 1.0, g.LineColor))

	start = renderer.Vector2{X: float32(-g.HalfWidth), Y: g.OriginY}
	end = renderer.Vector2{X: float32(g.HalfWidth), Y: g.OriginY}
	buffer.AddCommand(renderer.NewLineCommand(start, end, 1.0, g.LineColor))

	// Fade the grid line color for non-axes
	fadedColor := renderer.Color{
		R: uint8(float32(g.LineColor.R) * 0.25),
		G: uint8(float32(g.LineColor.G) * 0.25),
		B: uint8(float32(g.LineColor.B) * 0.25),
		A: g.LineColor.A,
	}

	// Draw vertical grid lines + labels
	for x := int(g.OriginX); x <= g.HalfWidth; x += g.Spacing {
		start := renderer.Vector2{X: float32(x), Y: float32(-g.HalfHeight)}
		end := renderer.Vector2{X: float32(x), Y: float32(g.HalfHeight)}
		buffer.AddCommand(renderer.NewLineCommand(start, end, 1.0, fadedColor))
		g.generateLabelCommand(fmt.Sprintf("%d", x), float32(x), g.OriginY, camera, buffer)
	}
	for x := int(g.OriginX); x >= -g.HalfWidth; x -= g.Spacing {
		start := renderer.Vector2{X: float32(x), Y: float32(-g.HalfHeight)}
		end := renderer.Vector2{X: float32(x), Y: float32(g.HalfHeight)}
		buffer.AddCommand(renderer.NewLineCommand(start, end, 1.0, fadedColor))
		g.generateLabelCommand(fmt.Sprintf("%d", x), float32(x), g.OriginY, camera, buffer)
	}

	// Draw horizontal grid lines + labels
	for y := int(g.OriginY); y <= g.HalfHeight; y += g.Spacing {
		start := renderer.Vector2{X: float32(-g.HalfWidth), Y: float32(y)}
		end := renderer.Vector2{X: float32(g.HalfWidth), Y: float32(y)}
		buffer.AddCommand(renderer.NewLineCommand(start, end, 1.0, fadedColor))
		g.generateLabelCommand(fmt.Sprintf("%d", y), g.OriginX, float32(y), camera, buffer)
	}
	for y := int(g.OriginY); y >= -g.HalfHeight; y -= g.Spacing {
		start := renderer.Vector2{X: float32(-g.HalfWidth), Y: float32(y)}
		end := renderer.Vector2{X: float32(g.HalfWidth), Y: float32(y)}
		buffer.AddCommand(renderer.NewLineCommand(start, end, 1.0, fadedColor))
		g.generateLabelCommand(fmt.Sprintf("%d", y), g.OriginX, float32(y), camera, buffer)
	}
}

// generateLabelCommand creates a label text command near the grid line in world space.
func (g *Grid) generateLabelCommand(text string, worldX, worldY float32, camera rl.Camera2D, buffer *renderer.CommandBuffer) {
	// Draw label slightly offset so it doesn't overlap the line
	offsetX, offsetY := float32(2), float32(2)
	position := renderer.Vector2{X: worldX + offsetX, Y: worldY + offsetY}
	buffer.AddCommand(renderer.NewTextCommand(text, position, g.LabelFontSize, g.LabelColor))
}

// Draw renders the grid using a renderer (for backwards compatibility)
func (g *Grid) Draw(camera rl.Camera2D) {
	buffer := renderer.NewCommandBuffer()
	g.GenerateDrawCommands(camera, buffer)

	raylibRenderer := renderer.NewRaylibRenderer()
	buffer.Execute(raylibRenderer)
}
