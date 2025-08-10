package musicfont

import (
	"encoding/json"
	"fmt"
	"os"
)

// FontEngravingDefaults holds engraving default values
type FontEngravingDefaults map[string]float64

// LoadFontEngravingDefaults loads engravingDefaults from a font metadata JSON file
func LoadFontEngravingDefaults(fontMetaPath string) (FontEngravingDefaults, error) {
	data, err := os.ReadFile(fontMetaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font metadata: %w", err)
	}
	var raw struct {
		EngravingDefaults map[string]float64 `json:"engravingDefaults"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse engraving defaults: %w", err)
	}
	return FontEngravingDefaults(raw.EngravingDefaults), nil
}
