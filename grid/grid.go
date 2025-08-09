package grid

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Grid struct {
	Spacing       int     // distance in pixels between lines in world space
	HalfWidth     int     // half width of grid in world units
	HalfHeight    int     // half height of grid in world units
	OriginX       float32 // grid origin X coordinate in world space
	OriginY       float32 // grid origin Y coordinate in world space
	LabelColor    rl.Color
	LineColor     rl.Color
	LabelFontSize int32
}

func New(spacing, halfWidth, halfHeight int, labelFontSize int32) *Grid {
	return &Grid{
		Spacing:       spacing,
		HalfWidth:     halfWidth,
		HalfHeight:    halfHeight,
		OriginX:       0,
		OriginY:       0,
		LabelColor:    rl.DarkGray,
		LineColor:     rl.LightGray,
		LabelFontSize: labelFontSize,
	}
}

func (g *Grid) Draw(camera rl.Camera2D) {
	// Draw central axes
	rl.DrawLine(int32(g.OriginX), -int32(g.HalfHeight), int32(g.OriginX), int32(g.HalfHeight), g.LineColor)
	rl.DrawLine(-int32(g.HalfWidth), int32(g.OriginY), int32(g.HalfWidth), int32(g.OriginY), g.LineColor)

	// Draw vertical grid lines + labels
	for x := int(g.OriginX); x <= g.HalfWidth; x += g.Spacing {
		rl.DrawLine(int32(x), -int32(g.HalfHeight), int32(x), int32(g.HalfHeight), rl.Fade(g.LineColor, 0.25))
		g.drawLabel(fmt.Sprintf("%d", x), int32(x), int32(g.OriginY), camera)
	}
	for x := int(g.OriginX); x >= -g.HalfWidth; x -= g.Spacing {
		rl.DrawLine(int32(x), -int32(g.HalfHeight), int32(x), int32(g.HalfHeight), rl.Fade(g.LineColor, 0.25))
		g.drawLabel(fmt.Sprintf("%d", x), int32(x), int32(g.OriginY), camera)
	}

	// Draw horizontal grid lines + labels
	for y := int(g.OriginY); y <= g.HalfHeight; y += g.Spacing {
		rl.DrawLine(-int32(g.HalfWidth), int32(y), int32(g.HalfWidth), int32(y), rl.Fade(g.LineColor, 0.25))
		g.drawLabel(fmt.Sprintf("%d", y), int32(g.OriginX), int32(y), camera)
	}
	for y := int(g.OriginY); y >= -g.HalfHeight; y -= g.Spacing {
		rl.DrawLine(-int32(g.HalfWidth), int32(y), int32(g.HalfWidth), int32(y), rl.Fade(g.LineColor, 0.25))
		g.drawLabel(fmt.Sprintf("%d", y), int32(g.OriginX), int32(y), camera)
	}
}

// drawLabel draws the label text near the grid line in world space.
func (g *Grid) drawLabel(text string, worldX, worldY int32, camera rl.Camera2D) {
	// Draw label slightly offset so it doesn't overlap the line
	offsetX, offsetY := int32(2), int32(2)
	rl.DrawText(text, worldX+offsetX, worldY+offsetY, g.LabelFontSize, g.LabelColor)
}
