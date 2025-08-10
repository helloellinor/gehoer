package musicfont

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// LoadFont loads a font with given glyph subset
func LoadMusicFont(path string, size int32, glyphs []rune) (rl.Font, error) {
	font := rl.LoadFontEx(path, size, glyphs)
	if font.Texture.ID == 0 {
		return rl.Font{}, fmt.Errorf("failed to load font")
	}
	return font, nil
}
