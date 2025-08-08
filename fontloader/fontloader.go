package fontloader

import (
	"fmt"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// LoadSmuflFont loads the font at fontPath with only the specified runes (glyphs),
// and returns the loaded rl.Font.
// It logs fatally if loading fails.
func LoadSmuflFont(fontPath string, fontSize int32, runes []rune) rl.Font {
	font := rl.LoadFontEx(fontPath, fontSize, runes, int32(len(runes)))
	if font.Texture.ID == 0 {
		log.Fatalf("Failed to load font texture for runes: %v", runes)
	}
	return font
}

// Safe wrapper that returns error instead of fatal log
func LoadSmuflFontSafe(fontPath string, fontSize int32, runes []rune) (rl.Font, error) {
	font := rl.LoadFontEx(fontPath, fontSize, runes, int32(len(runes)))
	if font.Texture.ID == 0 {
		return rl.Font{}, fmt.Errorf("failed to load font texture for runes: %v", runes)
	}
	return font, nil
}
