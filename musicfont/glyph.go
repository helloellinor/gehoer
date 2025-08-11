package musicfont

import (
	"fmt"
)

// Glyph represents a single glyph with combined metadata
type Glyph struct {
	Name        string
	Codepoint   rune
	Description string
	BBox        GlyphBBox             // bounding box info (from bbox.go)
	Anchors     map[string][2]float64 // font-specific anchors
}

func (mf *MusicFont) GetGlyph(name string) (*Glyph, bool) {
	g, ok := mf.GlyphMap[name]
	if !ok {
		fmt.Printf("GetGlyph: glyph %q NOT found\n", name)
	} else {
		fmt.Printf("GetGlyph: glyph %q found\n", name)
	}
	return g, ok
}
