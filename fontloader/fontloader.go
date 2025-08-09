package fontloader

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadSmuflFont(fontPath string, fontSize int32, runes []rune) (rl.Font, error) {
	font := rl.LoadFontEx(fontPath, fontSize, runes, int32(len(runes)))
	if font.Texture.ID == 0 {
		return rl.Font{}, fmt.Errorf("failed to load font texture for runes: %v", runes)
	}
	return font, nil
}
