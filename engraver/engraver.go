package engraver

import (
	"gehoer/music"
	"gehoer/musicfont"
	"gehoer/units"

	rl "github.com/gen2brain/raylib-go/raylib"
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

// DrawGlyph helper to draw a glyph by its name, with color
func (e *Engraver) DrawGlyph(glyphName string, x, y float32, color rl.Color) {
	glyph, ok := e.MusicFont.GetGlyph(glyphName)
	if !ok {
		// log or silently ignore missing glyph
		return
	}
	rl.DrawTextEx(e.MusicFont.Font, string(glyph.Codepoint), rl.NewVector2(x, y), e.fontSize, 0, color)
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

// Draw renders the entire score (simplified example)
func (e *Engraver) Draw(originX, originY float32) {
	x := originX
	y := originY

	for _, measure := range e.Score.Measures {
		staffLength := e.MeasureLengthPx(measure)
		e.DrawStaff(x, y, staffLength, rl.Black)
		// Draw clef if present
		if measure.Clef != "" {
			e.DrawClef(measure.Clef, 0, 0, rl.Black)
			x += 40 // arbitrary advance, use glyph bbox width ideally
		}

		// Draw measure elements (notes/rests/etc)
		for _, elem := range measure.Elements {
			switch el := elem.(type) {
			case *music.Note:
				e.DrawNote(el, x, y, rl.Black)
				x += 20 // advance x by some spacing (replace with glyph bbox width)
			default:
				e.DrawGlyph(el.GlyphName(), x, y, rl.Black)
				x += 20
			}
		}

		// Draw barline (not shown)
		x += units.StaffSpacesToPixels(5) // some margin after each measure

	}
}
