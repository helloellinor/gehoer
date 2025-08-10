package musicfont

import (
	"encoding/json"
	"fmt"
	"os"

	"gehoer/smufl"
)

type MusicFont struct {
	Metadata          *smufl.Metadata
	Anchors           Anchors
	EngravingDefaults map[string]float64
	BoundingBoxes     map[string]GlyphBBox
	GlyphMap          map[string]*Glyph // <--- Add this field
}

func LoadMetadata(smuflRepoPath, fontMetaJSONPath string) (*MusicFont, error) {
	metadata, err := smufl.LoadMetadata(smuflRepoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load SMuFL metadata: %w", err)
	}

	anchors, err := loadAnchors(fontMetaJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load font anchors: %w", err)
	}

	engravingDefaults, err := loadEngravingDefaults(fontMetaJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load engraving defaults: %w", err)
	}

	boundingBoxes, err := loadBoundingBoxes(fontMetaJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load bounding boxes: %w", err)
	}

	return &MusicFont{
		Metadata:          metadata,
		Anchors:           anchors,
		EngravingDefaults: engravingDefaults,
		BoundingBoxes:     boundingBoxes,
	}, nil
}

// loadAnchors loads glyph anchors from a font-specific metadata JSON file (internal helper)
func loadAnchors(jsonPath string) (Anchors, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font metadata: %w", err)
	}
	var raw struct {
		GlyphsWithAnchors map[string]map[string][2]float64 `json:"glyphsWithAnchors"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse font anchors: %w", err)
	}
	return Anchors(raw.GlyphsWithAnchors), nil
}

// loadEngravingDefaults loads engravingDefaults from a font-specific metadata JSON file (internal helper)
func loadEngravingDefaults(jsonPath string) (map[string]float64, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font metadata: %w", err)
	}
	var raw struct {
		EngravingDefaults map[string]float64 `json:"engravingDefaults"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse engraving defaults: %w", err)
	}
	return raw.EngravingDefaults, nil
}

// loadBoundingBoxes loads glyph bounding boxes from a font-specific metadata JSON file (internal helper)
func loadBoundingBoxes(jsonPath string) (map[string]GlyphBBox, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font metadata: %w", err)
	}
	var raw struct {
		GlyphBoundingBoxes map[string]GlyphBBox `json:"glyphBoundingBoxes"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse bounding boxes: %w", err)
	}
	return raw.GlyphBoundingBoxes, nil
}
