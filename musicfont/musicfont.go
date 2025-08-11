package musicfont

import (
	"encoding/json"
	"fmt"
	"os"

	"gehoer/smufl"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MusicFont struct {
	Metadata          *smufl.Metadata
	Font              rl.Font
	Anchors           Anchors
	EngravingDefaults *EngravingDefaults
	BoundingBoxes     map[string]GlyphBBox
	GlyphMap          map[string]*Glyph
}

// EngravingDefaults represents the engravingDefaults JSON object with correct types.
type EngravingDefaults struct {
	ArrowShaftThickness        float64  `json:"arrowShaftThickness"`
	BarlineSeparation          float64  `json:"barlineSeparation"`
	BeamSpacing                float64  `json:"beamSpacing"`
	BeamThickness              float64  `json:"beamThickness"`
	BracketThickness           float64  `json:"bracketThickness"`
	DashedBarlineDashLength    float64  `json:"dashedBarlineDashLength"`
	DashedBarlineGapLength     float64  `json:"dashedBarlineGapLength"`
	DashedBarlineThickness     float64  `json:"dashedBarlineThickness"`
	HBarThickness              float64  `json:"hBarThickness"`
	HairpinThickness           float64  `json:"hairpinThickness"`
	LegerLineExtension         float64  `json:"legerLineExtension"`
	LegerLineThickness         float64  `json:"legerLineThickness"`
	LyricLineThickness         float64  `json:"lyricLineThickness"`
	OctaveLineThickness        float64  `json:"octaveLineThickness"`
	PedalLineThickness         float64  `json:"pedalLineThickness"`
	RepeatBarlineDotSeparation float64  `json:"repeatBarlineDotSeparation"`
	RepeatEndingLineThickness  float64  `json:"repeatEndingLineThickness"`
	SlurEndpointThickness      float64  `json:"slurEndpointThickness"`
	SlurMidpointThickness      float64  `json:"slurMidpointThickness"`
	StaffLineThickness         float64  `json:"staffLineThickness"`
	StemThickness              float64  `json:"stemThickness"`
	SubBracketThickness        float64  `json:"subBracketThickness"`
	TextEnclosureThickness     float64  `json:"textEnclosureThickness"`
	TextFontFamily             []string `json:"textFontFamily"`
	ThickBarlineThickness      float64  `json:"thickBarlineThickness"`
	ThinBarlineThickness       float64  `json:"thinBarlineThickness"`
	ThinThickBarlineSeparation float64  `json:"thinThickBarlineSeparation"`
	TieEndpointThickness       float64  `json:"tieEndpointThickness"`
	TieMidpointThickness       float64  `json:"tieMidpointThickness"`
	TupletBracketThickness     float64  `json:"tupletBracketThickness"`
}

func LoadMusicFont(smuflRepoPath, fontMetaJSONPath, fontFilePath string, fontSize int32) (*MusicFont, error) {
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

	// Collect all needed runes from metadata:
	var runes []rune
	for name := range metadata.Glyphs {
		r, err := metadata.GetGlyphRune(name)
		if err == nil {
			runes = append(runes, r)
		}
	}

	// Load font with subset of runes (like old fontloader)
	font := rl.LoadFontEx(fontFilePath, int32(units.FontLoadSizePx), runes, int32(len(runes)))
	if font.Texture.ID == 0 {
		return nil, fmt.Errorf("failed to load font texture")
	}

	mf := &MusicFont{
		Metadata:          metadata,
		Font:              font,
		Anchors:           anchors,
		EngravingDefaults: engravingDefaults,
		BoundingBoxes:     boundingBoxes,
	}

	mf.BuildGlyphMap()

	return mf, nil

}

// BuildGlyphMap builds enriched glyph structs combining SMuFL metadata and font-specific data
func (mf *MusicFont) BuildGlyphMap() {
	mf.GlyphMap = make(map[string]*Glyph)
	for name, smGlyph := range mf.Metadata.Glyphs {
		r, err := mf.Metadata.GetGlyphRune(name)
		if err != nil {
			// skip glyphs with invalid rune
			continue
		}

		bbox := mf.BoundingBoxes[name]
		anchors := mf.Anchors[name]

		mf.GlyphMap[name] = &Glyph{
			Name:        name,
			Codepoint:   r,
			Description: smGlyph.Description,
			BBox:        bbox,
			Anchors:     anchors,
		}
	}
}

// --- helpers ---

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

// loadEngravingDefaults loads engraving defaults from a JSON file.
func loadEngravingDefaults(jsonPath string) (*EngravingDefaults, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font metadata: %w", err)
	}

	var raw struct {
		EngravingDefaults EngravingDefaults `json:"engravingDefaults"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse engraving defaults: %w", err)
	}

	return &raw.EngravingDefaults, nil
}

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
