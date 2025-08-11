package engraver

import (
	"gehoer/music"
	"gehoer/musicfont"
	"gehoer/renderer"
	"gehoer/units"
)

type Engraver struct {
	Score     *music.Score
	MusicFont *musicfont.MusicFont
	fontSize  float32
}

func NewEngraver(score *music.Score, musicFont *musicfont.MusicFont) *Engraver {
	return &Engraver{
		Score:     score,
		MusicFont: musicFont,
		fontSize:  units.FontRenderSizePx, // load from your global settings
	}
}

// CreateGlyphCommand helper to create a glyph command by its name, with color
func (e *Engraver) CreateGlyphCommand(glyphName string, x, y float32, color renderer.Color) *renderer.GlyphCommand {
	glyph, ok := e.MusicFont.GetGlyph(glyphName)
	if !ok {
		// Return nil for missing glyph - caller should check
		return nil
	}
	cmd := renderer.NewGlyphCommand(e.MusicFont.Font, glyph.Codepoint, renderer.Vector2{X: x, Y: y}, e.fontSize, color)
	return &cmd
}

func (e *Engraver) MeasureLengthPx(measure *music.Measure) float32 {
	totalWidth := float32(0)
	spacing := units.StaffSpacesToPixels(2) // pixels between glyphs, tweak as needed

	for _, elem := range measure.Elements {
		glyphName := elem.GlyphName()
		glyph, ok := e.MusicFont.GetGlyph(glyphName)
		if !ok {
			// fallback width for missing glyphs (e.g. 10 px)
			totalWidth += units.StaffSpacesToPixels(2) + spacing
			continue
		}
		bboxWidth := e.bboxWidthInPixels(glyph.BBox)
		totalWidth += bboxWidth + spacing
	}
	return totalWidth
}

// Helper to calculate pixel width from GlyphBBox
func (e *Engraver) bboxWidthInPixels(bbox musicfont.GlyphBBox) float32 {
	// bbox.NE[0] = northeast x; bbox.SW[0] = southwest x
	widthStaffSpaces := bbox.NE[0] - bbox.SW[0]
	return units.StaffSpacesToPixels(float32(widthStaffSpaces))
}

// GenerateDrawCommands creates draw commands for the entire score (simplified example)
func (e *Engraver) GenerateDrawCommands(originX, originY float32, buffer *renderer.CommandBuffer) {
	x := originX
	y := originY

	for _, measure := range e.Score.Measures {
		staffLength := e.MeasureLengthPx(measure)
		e.GenerateStaffCommands(x, y, staffLength, renderer.Black, buffer)

		// Draw clef if present
		if measure.Clef != "" {
			if cmd := e.CreateGlyphCommand(measure.Clef, x, y, renderer.Black); cmd != nil {
				buffer.AddCommand(*cmd)
			}
			x += 40 // arbitrary advance, use glyph bbox width ideally
		}

		// Draw measure elements (notes/rests/etc)
		for _, elem := range measure.Elements {
			switch el := elem.(type) {
			case *music.Note:
				e.GenerateNoteCommands(el, x, y, renderer.Black, buffer)
				x += 20 // advance x by some spacing (replace with glyph bbox width)
			default:
				if cmd := e.CreateGlyphCommand(el.GlyphName(), x, y, renderer.Black); cmd != nil {
					buffer.AddCommand(*cmd)
				}
				x += 20
			}
		}

		// Draw barline (not shown)
		x += units.StaffSpacesToPixels(5) // some margin after each measure
	}
}
