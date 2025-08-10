package musicfont

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Anchor represents a named anchor point with coordinates
type Anchor struct {
	Name string
	X    float64
	Y    float64
}

// Anchors maps glyph names to anchor points: anchorName -> [x, y] in staff spaces
type Anchors map[string]map[string][2]float64

// GlyphAnchors holds all anchors for a single glyph as a map from anchor name to Anchor
//type GlyphAnchors map[string]Anchor

// Distance calculates Euclidean distance ignoring the name
func (a Anchor) Distance(b Anchor) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// LoadAnchors loads glyph anchors from a JSON file.
func LoadAnchors(path string) (Anchors, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read anchors JSON: %w", err)
	}

	var raw struct {
		GlyphsWithAnchors Anchors `json:"glyphsWithAnchors"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal anchors JSON: %w", err)
	}

	return raw.GlyphsWithAnchors, nil
}
