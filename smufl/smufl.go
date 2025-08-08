package smufl

import (
	"encoding/json"
	"fmt"
	"os"
)

type SmuflData struct {
	GlyphBBoxes map[string]GlyphBBox `json:"glyphBBoxes"`
}

type GlyphBBox struct {
	BBoxNE []float64 `json:"bBoxNE"` // [x, y]
	BBoxSW []float64 `json:"bBoxSW"` // [x, y]
}

// LoadSmuflMetadata loads SMuFL JSON metadata from file
func LoadSmuflMetadata(path string) (*SmuflData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read SMuFL metadata: %w", err)
	}

	var smufl SmuflData
	if err := json.Unmarshal(data, &smufl); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SMuFL metadata: %w", err)
	}

	return &smufl, nil
}
