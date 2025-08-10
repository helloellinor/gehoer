package musicfont

import (
	"encoding/json"
	"fmt"
	"os"
)

// GlyphBBox defines the bounding box for a glyph in staff spaces.
type GlyphBBox struct {
	SW [2]float64 `json:"sw"` // southwest corner
	NE [2]float64 `json:"ne"` // northeast corner
}

// LoadBoundingBoxes loads glyph bounding boxes from a JSON file.
func LoadBoundingBoxes(path string) (map[string]GlyphBBox, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read bounding boxes JSON: %w", err)
	}

	var bboxes map[string]GlyphBBox
	if err := json.Unmarshal(data, &bboxes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bounding boxes JSON: %w", err)
	}

	return bboxes, nil
}
