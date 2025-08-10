package musicfont

// Glyph represents a single glyph with combined metadata
type Glyph struct {
	Name        string
	Codepoint   rune
	Description string
	BBox        GlyphBBox             // bounding box info (from bbox.go)
	Anchors     map[string][2]float64 // font-specific anchors
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

// GetGlyph returns the enriched Glyph by name, if present
func (mf *MusicFont) GetGlyph(name string) (*Glyph, bool) {
	g, ok := mf.GlyphMap[name]
	return g, ok
}
