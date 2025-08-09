package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// SMuFLGlyph represents a single SMuFL glyph definition
type SMuFLGlyph struct {
	Codepoint           string `json:"codepoint"`
	AlternateCodepoint  string `json:"alternateCodepoint,omitempty"`
	Description         string `json:"description"`
}

// SMuFLRange represents a range of related glyphs
type SMuFLRange struct {
	Description string   `json:"description"`
	Glyphs      []string `json:"glyphs"`
	RangeStart  string   `json:"range_start"`
	RangeEnd    string   `json:"range_end"`
}

// SMuFLClass represents a classification of glyphs (from classes.json)
// Note: classes.json actually stores arrays directly, not objects with description/keys
type SMuFLClass []string

// SMuFLMetadata holds all SMuFL metadata
type SMuFLMetadata struct {
	Glyphs  map[string]SMuFLGlyph `json:"glyphs"`
	Ranges  map[string]SMuFLRange `json:"ranges"`
	Classes map[string]SMuFLClass `json:"classes"`
}

// FontAnchors holds glyph attachment points from a font's metadata (glyphsWithAnchors)
// Map: glyphName -> anchorName -> [x,y] in staff spaces
type FontAnchors map[string]map[string][2]float64

// FontEngravingDefaults holds useful engraving defaults from a font's metadata
// Values are in staff spaces unless noted otherwise
type FontEngravingDefaults map[string]float64

// GlyphBBox represents a glyph bounding box (compatible with old format)
type GlyphBBox struct {
	SW [2]float64 `json:"sw"` // Southwest corner
	NE [2]float64 `json:"ne"` // Northeast corner
}

// LoadSMuFLMetadata loads all SMuFL metadata from the official repository
func LoadSMuFLMetadata(smuflRepoPath string) (*SMuFLMetadata, error) {
	metadata := &SMuFLMetadata{
		Glyphs:  make(map[string]SMuFLGlyph),
		Ranges:  make(map[string]SMuFLRange),
		Classes: make(map[string]SMuFLClass),
	}

	// Load glyphnames.json
	glyphNamesPath := filepath.Join(smuflRepoPath, "metadata", "glyphnames.json")
	if err := loadJSONFile(glyphNamesPath, &metadata.Glyphs); err != nil {
		return nil, fmt.Errorf("failed to load glyphnames.json: %w", err)
	}

	// Load ranges.json
	rangesPath := filepath.Join(smuflRepoPath, "metadata", "ranges.json")
	if err := loadJSONFile(rangesPath, &metadata.Ranges); err != nil {
		return nil, fmt.Errorf("failed to load ranges.json: %w", err)
	}

	// Load classes.json
	classesPath := filepath.Join(smuflRepoPath, "metadata", "classes.json")
	if err := loadJSONFile(classesPath, &metadata.Classes); err != nil {
		return nil, fmt.Errorf("failed to load classes.json: %w", err)
	}

	return metadata, nil
}

// loadJSONFile loads a JSON file into the provided interface
func loadJSONFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// GetGlyphRune returns the rune for a glyph name
func (s *SMuFLMetadata) GetGlyphRune(glyphName string) (rune, error) {
	glyph, exists := s.Glyphs[glyphName]
	if !exists {
		return 0, fmt.Errorf("glyph %s not found", glyphName)
	}
	
	// Parse the Unicode codepoint (format: "U+E123")
	codeStr := strings.TrimPrefix(glyph.Codepoint, "U+")
	code, err := strconv.ParseInt(codeStr, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse codepoint %s: %w", glyph.Codepoint, err)
	}
	
	return rune(code), nil
}

// GetGlyphsByRange returns all glyphs in a specific range
func (s *SMuFLMetadata) GetGlyphsByRange(rangeName string) ([]string, error) {
	rangeInfo, exists := s.Ranges[rangeName]
	if !exists {
		return nil, fmt.Errorf("range %s not found", rangeName)
	}
	return rangeInfo.Glyphs, nil
}

// GetGlyphsByClass returns all glyphs in a specific class
func (s *SMuFLMetadata) GetGlyphsByClass(className string) ([]string, error) {
	classInfo, exists := s.Classes[className]
	if !exists {
		return nil, fmt.Errorf("class %s not found", className)
	}
	return []string(classInfo), nil
}

// GetAllNoteheads returns all notehead glyphs
func (s *SMuFLMetadata) GetAllNoteheads() ([]string, error) {
	noteheads := []string{}
	
	// Get basic noteheads
	if basicNoteheads, err := s.GetGlyphsByRange("noteheads"); err == nil {
		noteheads = append(noteheads, basicNoteheads...)
	}
	
	// Get supplemental noteheads
	if suppNoteheads, err := s.GetGlyphsByRange("noteheadsSupplement"); err == nil {
		noteheads = append(noteheads, suppNoteheads...)
	}
	
	// Get shape note noteheads
	if shapeNoteheads, err := s.GetGlyphsByRange("shapeNoteHeads"); err == nil {
		noteheads = append(noteheads, shapeNoteheads...)
	}
	
	// Get shape note noteheads supplement
	if shapeNoteSuppNoteheads, err := s.GetGlyphsByRange("shapeNoteHeadsSupplement"); err == nil {
		noteheads = append(noteheads, shapeNoteSuppNoteheads...)
	}
	
	// Remove duplicates and sort
	seen := make(map[string]bool)
	uniqueNoteheads := []string{}
	for _, notehead := range noteheads {
		if !seen[notehead] {
			seen[notehead] = true
			uniqueNoteheads = append(uniqueNoteheads, notehead)
		}
	}
	sort.Strings(uniqueNoteheads)
	
	return uniqueNoteheads, nil
}

// GetBasicNoteGlyphs returns the standard note glyphs (whole, half, quarter, etc.)
func (s *SMuFLMetadata) GetBasicNoteGlyphs() (map[string]rune, error) {
	basicNotes := map[string]string{
		"noteDoubleWhole":    "noteDoubleWhole",
		"noteWhole":          "noteWhole", 
		"noteHalfUp":         "noteHalfUp",
		"noteHalfDown":       "noteHalfDown",
		"noteQuarterUp":      "noteQuarterUp",
		"noteQuarterDown":    "noteQuarterDown",
		"note8thUp":          "note8thUp",
		"note8thDown":        "note8thDown",
		"note16thUp":         "note16thUp",
		"note16thDown":       "note16thDown",
		"note32ndUp":         "note32ndUp",
		"note32ndDown":       "note32ndDown",
		"note64thUp":         "note64thUp",
		"note64thDown":       "note64thDown",
		"note128thUp":        "note128thUp",
		"note128thDown":      "note128thDown",
		"note256thUp":        "note256thUp",
		"note256thDown":      "note256thDown",
		"note512thUp":        "note512thUp",
		"note512thDown":      "note512thDown",
		"note1024thUp":       "note1024thUp",
		"note1024thDown":     "note1024thDown",
	}

	result := make(map[string]rune)
	for key, glyphName := range basicNotes {
		if r, err := s.GetGlyphRune(glyphName); err == nil {
			result[key] = r
		}
		// If specific glyph not found, try fallback patterns
	}
	
	return result, nil
}

// GetBasicRestGlyphs returns the standard rest glyphs
func (s *SMuFLMetadata) GetBasicRestGlyphs() (map[string]rune, error) {
	basicRests := map[string]string{
		"restMaxima":         "restMaxima",
		"restLonga":          "restLonga",
		"restDoubleWhole":    "restDoubleWhole",
		"restWhole":          "restWhole",
		"restHalf":           "restHalf",
		"restQuarter":        "restQuarter",
		"rest8th":            "rest8th",
		"rest16th":           "rest16th",
		"rest32nd":           "rest32nd",
		"rest64th":           "rest64th",
		"rest128th":          "rest128th",
		"rest256th":          "rest256th",
		"rest512th":          "rest512th",
		"rest1024th":         "rest1024th",
	}

	result := make(map[string]rune)
	for key, glyphName := range basicRests {
		if r, err := s.GetGlyphRune(glyphName); err == nil {
			result[key] = r
		}
	}
	
	return result, nil
}

// GetClefGlyphs returns common clef glyphs
func (s *SMuFLMetadata) GetClefGlyphs() (map[string]rune, error) {
	clefs := map[string]string{
		"gClef":              "gClef",
		"fClef":              "fClef", 
		"cClef":              "cClef",
		"percClef":           "percClef",
		"gClefOttavaAlta":    "gClefOttavaAlta",
		"gClefOttavaBassa":   "gClefOttavaBassa",
		"fClefOttavaAlta":    "fClefOttavaAlta", 
		"fClefOttavaBassa":   "fClefOttavaBassa",
		"cClefAlto":          "cClefAlto",
		"6stringTabClef":     "6stringTabClef",
		"4stringTabClef":     "4stringTabClef",
	}

	result := make(map[string]rune)
	for key, glyphName := range clefs {
		if r, err := s.GetGlyphRune(glyphName); err == nil {
			result[key] = r
		}
	}
	
	return result, nil
}

// GetAccidentalGlyphs returns accidental glyphs
func (s *SMuFLMetadata) GetAccidentalGlyphs() (map[string]rune, error) {
	accidentals := map[string]string{
		"accidentalFlat":          "accidentalFlat",
		"accidentalNatural":       "accidentalNatural",
		"accidentalSharp":         "accidentalSharp",
		"accidentalDoubleFlat":    "accidentalDoubleFlat",
		"accidentalDoubleSharp":   "accidentalDoubleSharp",
		"accidentalTripleFlat":    "accidentalTripleFlat",
		"accidentalTripleSharp":   "accidentalTripleSharp",
		"accidentalQuarterToneSharpStein": "accidentalQuarterToneSharpStein",
		"accidentalQuarterToneFlatStein":  "accidentalQuarterToneFlatStein",
	}

	result := make(map[string]rune)
	for key, glyphName := range accidentals {
		if r, err := s.GetGlyphRune(glyphName); err == nil {
			result[key] = r
		}
	}
	
	return result, nil
}

// GetAllGlyphRunes returns runes for all specified glyph names
func (s *SMuFLMetadata) GetAllGlyphRunes(glyphNames []string) ([]rune, error) {
	runes := make([]rune, 0, len(glyphNames))
	
	for _, name := range glyphNames {
		if r, err := s.GetGlyphRune(name); err == nil {
			runes = append(runes, r)
		} else {
			// Log missing glyph but continue
			fmt.Printf("Warning: glyph %s not found: %v\n", name, err)
		}
	}
	
	return runes, nil
}

// ListAvailableRanges returns all available range names
func (s *SMuFLMetadata) ListAvailableRanges() []string {
	ranges := make([]string, 0, len(s.Ranges))
	for name := range s.Ranges {
		ranges = append(ranges, name)
	}
	sort.Strings(ranges)
return ranges
}

// LoadFontAnchors loads glyphsWithAnchors from a font-specific metadata JSON file
func LoadFontAnchors(fontMetaPath string) (FontAnchors, error) {
	data, err := os.ReadFile(fontMetaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font metadata: %w", err)
	}
	// Minimal struct to unmarshal only the glyphsWithAnchors section
	var raw struct {
		GlyphsWithAnchors map[string]map[string][2]float64 `json:"glyphsWithAnchors"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse font anchors: %w", err)
	}
return FontAnchors(raw.GlyphsWithAnchors), nil
}

// LoadFontEngravingDefaults loads engravingDefaults from a font-specific metadata JSON file
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

// ListAvailableClasses returns all available class names  
func (s *SMuFLMetadata) ListAvailableClasses() []string {
	classes := make([]string, 0, len(s.Classes))
	for name := range s.Classes {
		classes = append(classes, name)
	}
	sort.Strings(classes)
	return classes
}

// CreateGlyphBoundingBoxMap creates a map compatible with the old format
// Note: The official SMuFL repo doesn't include bounding box data in the main metadata files
// This would need to be loaded separately from font-specific files
func (s *SMuFLMetadata) CreateGlyphBoundingBoxMap() map[string]GlyphBBox {
	// This is a placeholder - actual bounding box data would come from 
	// font-specific JSON files that are not part of the core SMuFL spec
	bboxes := make(map[string]GlyphBBox)
	
	// Add some reasonable default bounding boxes for common glyphs
	// These are estimates and would need to be replaced with actual font data
	commonGlyphs := map[string]GlyphBBox{
		"noteWhole": {
			SW: [2]float64{-0.75, -0.5},
			NE: [2]float64{0.75, 0.5},
		},
		"noteHalfUp": {
			SW: [2]float64{-0.6, -0.5},
			NE: [2]float64{0.6, 2.5},
		},
		"noteQuarterUp": {
			SW: [2]float64{-0.6, -0.5},
			NE: [2]float64{0.6, 2.5},
		},
		"note8thUp": {
			SW: [2]float64{-0.6, -0.5},
			NE: [2]float64{0.8, 2.5},
		},
	}
	
	// Copy common glyphs to the map
	for name, bbox := range commonGlyphs {
		bboxes[name] = bbox
	}
	
	return bboxes
}
