package engraver

import (
	"gehoer/smufl"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Clef struct{}

func (c *Clef) Draw(clefName string, originX, originY float32, font rl.Font, smufl *smufl.Metadata) {
	clefRune, err := smufl.GetGlyphRune(clefName)
	if err != nil {
		return // silently ignore missing glyph
	}

	DrawGlyph(font, clefRune, originX+units.StaffSpacesToPixels(0.5), originY, 0, rl.Black)
}
