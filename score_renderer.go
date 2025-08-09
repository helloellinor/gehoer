package main

import (
	"gehoer/metadata"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ScoreRenderer handles rendering of musical scores
type ScoreRenderer struct {
	smuflMetadata *metadata.SMuFLMetadata
	font          rl.Font
	bboxMap       map[string]metadata.GlyphBBox
	anchors       metadata.FontAnchors
	defaults      metadata.FontEngravingDefaults

	// playback
	playing     bool
	startTime   float64
	tempoBPM    int
	schedule    []scheduledEvent
	activeIndex int
	stream      rl.AudioStream
	audioReady  bool
	sampleRate  int
	samplePhase float64

	// rendering index
	elementCounter int
}

type scheduledEvent struct {
	start  float64
	end    float64
	isNote bool
	pitch  int // MIDI
}

// NewScoreRenderer creates a new score renderer
func NewScoreRenderer(smuflMetadata *metadata.SMuFLMetadata, font rl.Font, bboxMap map[string]metadata.GlyphBBox, anchors metadata.FontAnchors, defaults metadata.FontEngravingDefaults) *ScoreRenderer {
	return &ScoreRenderer{
		smuflMetadata: smuflMetadata,
		font:          font,
		bboxMap:       bboxMap,
		anchors:       anchors,
		defaults:      defaults,
		// audio defaults
		sampleRate: 44100,
	}
}

// RenderScore renders an entire score at the specified position
func (sr *ScoreRenderer) RenderScore(score *Score, startX, startY float32) {
	// initialize audio stream lazily
	if !sr.audioReady {
		rl.InitAudioDevice()
		sr.stream = rl.LoadAudioStream(uint32(sr.sampleRate), 32, 1)
		rl.PlayAudioStream(sr.stream)
		sr.audioReady = true
		sr.tempoBPM = score.Tempo
		sr.BuildSchedule(score)
		sr.StartPlayback()
	}
	// update active note based on time
	sr.updatePlayback()
	// reset element counter each frame
	sr.elementCounter = 0
	currentX := startX
	currentY := startY

	// Draw title (move up by 100px)
	titleText := score.Title + " (" + score.KeySignature.String() + ", " + score.TimeSignature.String() + ")"
	rl.DrawText(titleText, int32(startX), int32(startY-150), 24, rl.Black)

	// Calculate measure width based on time signature and content
	measureWidth := StaffSpacesToPixels(20) // 20 staff spaces per measure
	staffHeight := StaffSpacesToPixels(8)   // 8 staff spaces high (4 above + 4 below staff)

	measuresPerLine := 4 // Number of measures per staff line
	DrawClef("treble", startX, startY-StaffSpacePx*9, sr.font, sr.smuflMetadata, sr.bboxMap)
	currentX += 4 * StaffSpacePx // Adjust startX to account for clef width
	for i, measure := range score.Measures {
		// Start new line every measuresPerLine measures
		if i%measuresPerLine == 0 && i > 0 {
			currentX = startX
			currentY += staffHeight * 2 // Space between staff lines

		}

		// Draw measure
		sr.RenderMeasure(measure, currentX, currentY, measureWidth)

		currentX += measureWidth
	}
}

// RenderMeasure renders a single measure at the specified position
func (sr *ScoreRenderer) RenderMeasure(measure *Measure, startX, startY, width float32) {
	// Draw staff lines for this measure
	DrawStaff(startX, startY, width/StaffSpacesToPixels(1)) // Convert width to staff spaces
	// Draw measure number
	measureNumText := string(rune(measure.Number + '0'))
	if measure.Number >= 10 {
		measureNumText = "1" + string(rune((measure.Number%10)+'0'))
	}
	rl.DrawText(measureNumText, int32(startX), int32(startY-StaffSpacesToPixels(3)), 16, rl.DarkGray)

	// Calculate note positions within the measure
	elementCount := len(measure.Elements)
	if elementCount == 0 {
		return
	}

	// Proportional spacing by duration using measure logic
	leftMargin := StaffSpacesToPixels(1)
	rightMargin := StaffSpacesToPixels(1)
	positions := measure.ElementPositions(width, leftMargin, rightMargin)
	for i, element := range measure.Elements {
		x := startX + positions[i]
		switch e := element.(type) {
		case Note:
			sr.RenderNote(e, x, startY)
		case Rest:
			sr.RenderRest(e, x, startY)
		}
		sr.elementCounter++
	}

	// Draw measure bar at the end - exactly 1em tall from bottom to top staff line
	DrawMeasureBar(startX+width, startY)
}

// RenderNote renders a single note at the specified position
func (sr *ScoreRenderer) RenderNote(note Note, x, y float32) {
	// Calculate vertical position based on staff line
	// staffLine 0 = bottom line of staff, 1 = second line, etc.
	// Negative staffLine = below staff, positive = above staff when > 4
	noteY := y - StaffSpacesToPixels(float32(note.StaffLine)*0.5) // Each staff line is 0.5 staff spaces apart

	// Choose notehead by duration
	var noteheadName string
	switch note.Duration {
	case WholeNote:
		noteheadName = "noteheadWhole"
	case HalfNote:
		noteheadName = "noteheadHalf"
	default:
		noteheadName = "noteheadBlack"
	}
	// Resolve rune
	noteheadRune, err := sr.smuflMetadata.GetGlyphRune(noteheadName)
	if err != nil {
		// Fallback to previously hardcoded rune
		noteheadRune = note.Duration.GetSMUFLRune()
	}
	// choose color if active
	color := rl.Black
	if sr.elementCounter == sr.activeIndex {
		color = rl.Red
	}
	// Draw the notehead
	if bbox, hasBBox := sr.bboxMap[noteheadName]; hasBBox {
		drawGlyphWithBBox(sr.font, noteheadRune, bbox, x, noteY, 0)
	} else {
		rl.DrawTextEx(sr.font, string(noteheadRune), rl.NewVector2(x, noteY-(FontRenderSize/2)), FontRenderSize, 0.0, color)
	}

	// Draw stem and flag for notes shorter than whole notes
	if note.Duration != WholeNote {
		// Determine stem direction: below middle line -> up; on/above -> down
		stemUp := note.StaffLine < 3
		stemLength := StaffSpacesToPixels(3.5)
		var stemX float32
		var stemStartY, stemEndY float32
		stemThickness := float32(2.0)
		if sr.defaults != nil {
			if v, ok := sr.defaults["stemThickness"]; ok {
				stemThickness = StaffSpacesToPixels(float32(v))
			}
		}

		// Prefer font anchors for exact attachment; fallback to bbox/heuristic
		if sr.anchors != nil {
			if stemUp {
				if a, ok := sr.anchors[noteheadName]["stemUpSE"]; ok {
					stemX = x + StaffSpacesToPixels(float32(a[0]))
					stemStartY = noteY - StaffSpacesToPixels(float32(a[1]))
					stemEndY = stemStartY - stemLength
				}
			} else {
				if a, ok := sr.anchors[noteheadName]["stemDownNW"]; ok {
					stemX = x + StaffSpacesToPixels(float32(a[0]))
					stemStartY = noteY - StaffSpacesToPixels(float32(a[1]))
					stemEndY = stemStartY + stemLength
				}
			}
		}
		// If anchors didn't set stemX, fallback to bbox
		if stemX == 0 {
			if bbox, ok := sr.bboxMap[noteheadName]; ok {
				if stemUp {
					stemX = x + StaffSpacesToPixels(float32(bbox.NE[0]))
					stemStartY = noteY
					stemEndY = noteY - stemLength
				} else {
					stemX = x + StaffSpacesToPixels(float32(bbox.SW[0]))
					stemStartY = noteY
					stemEndY = noteY + stemLength
				}
			} else {
				if stemUp {
					stemX = x + StaffSpacesToPixels(0.5)
					stemStartY = noteY
					stemEndY = noteY - stemLength
				} else {
					stemX = x - StaffSpacesToPixels(0.5)
					stemStartY = noteY
					stemEndY = noteY + stemLength
				}
			}
		}

		// Draw stem (offset inward by half thickness to meet notehead edge)
		var stemDrawX float32 = stemX
		if stemUp {
			stemDrawX = stemX - stemThickness/2
		} else {
			stemDrawX = stemX + stemThickness/2
		}
		rl.DrawLineEx(
			rl.NewVector2(stemDrawX, stemStartY),
			rl.NewVector2(stemDrawX, stemEndY),
			stemThickness, color)
		// Draw flag if needed
		flagName := sr.flagGlyphName(note.Duration, stemUp)
		if flagName != "" {
			if flagRune, err := sr.smuflMetadata.GetGlyphRune(flagName); err == nil {
				// Use flag anchors to align anchor point to stem end
				dx := float32(0)
				dy := float32(0)
				if sr.anchors != nil {
					if stemUp {
						if a, ok := sr.anchors[flagName]["stemUpNW"]; ok {
							dx = StaffSpacesToPixels(float32(a[0]))
							dy = StaffSpacesToPixels(float32(a[1]))
						}
					} else {
						if a, ok := sr.anchors[flagName]["stemDownSW"]; ok {
							dx = StaffSpacesToPixels(float32(a[0]))
							dy = StaffSpacesToPixels(float32(a[1]))
						}
					}
				}
				flagPos := rl.NewVector2(stemDrawX-dx, stemEndY+dy)
				rl.DrawTextEx(sr.font, string(flagRune), flagPos, FontRenderSize, 0.0, color)
			}
		}
	}

	// Draw accidental if present
	if note.Accidental != "" {
		accidentalX := x - StaffSpacesToPixels(1.5) // Position accidental to the left of the note
		sr.RenderAccidental(note.Accidental, accidentalX, noteY)
	}

	// Add ledger lines if needed (notes above or below the staff)
	if note.StaffLine < 0 {
		// Draw ledger lines below staff; center horizontally on notehead
		centerX := x
		ledgerY := y + StaffSpacesToPixels(0.5) // First ledger line below staff
		if bbox, ok := sr.bboxMap[noteheadName]; ok {
			centerX = x + StaffSpacesToPixels(float32((bbox.SW[0]+bbox.NE[0])/2))
		}
		thickness := float32(1.0)
		if sr.defaults != nil {
			if v, ok := sr.defaults["staffLineThickness"]; ok {
				thickness = StaffSpacesToPixels(float32(v))
			}
		}
		for line := 0; line > note.StaffLine; line -= 2 {
			if line <= 0 {
				rl.DrawLineEx(
					rl.NewVector2(centerX-StaffSpacesToPixels(0.5), ledgerY-StaffSpacesToPixels(float32(line)*0.5)),
					rl.NewVector2(centerX+StaffSpacesToPixels(0.5), ledgerY-StaffSpacesToPixels(float32(line)*0.5)),
					thickness, rl.Black)
			}
		}
	} else if note.StaffLine > 8 {
		// Draw ledger lines above staff; center horizontally on notehead
		centerX := x
		ledgerY := y - StaffSpacesToPixels(4.5) // First ledger line above staff
		if bbox, ok := sr.bboxMap[noteheadName]; ok {
			centerX = x + StaffSpacesToPixels(float32((bbox.SW[0]+bbox.NE[0])/2))
		}
		thickness := float32(1.0)
		if sr.defaults != nil {
			if v, ok := sr.defaults["staffLineThickness"]; ok {
				thickness = StaffSpacesToPixels(float32(v))
			}
		}
		for line := 10; line <= note.StaffLine; line += 2 {
			rl.DrawLineEx(
				rl.NewVector2(centerX-StaffSpacesToPixels(0.75), ledgerY-StaffSpacesToPixels(float32(line-10)*0.5)),
				rl.NewVector2(centerX+StaffSpacesToPixels(0.75), ledgerY-StaffSpacesToPixels(float32(line-10)*0.5)),
				thickness, rl.Black)
		}
	}
}

// flagGlyphName returns the correct SMuFL flag glyph name for a duration and stem direction
func (sr *ScoreRenderer) flagGlyphName(nv NoteValue, stemUp bool) string {
	if stemUp {
		switch nv {
		case EighthNote:
			return "flag8thUp"
		case SixteenthNote:
			return "flag16thUp"
		case ThirtySecondNote:
			return "flag32ndUp"
		case SixtyFourthNote:
			return "flag64thUp"
		default:
			return ""
		}
	} else {
		switch nv {
		case EighthNote:
			return "flag8thDown"
		case SixteenthNote:
			return "flag16thDown"
		case ThirtySecondNote:
			return "flag32ndDown"
		case SixtyFourthNote:
			return "flag64thDown"
		default:
			return ""
		}
	}
}

// RenderRest renders a musical rest at the specified position
func (sr *ScoreRenderer) RenderRest(rest Rest, x, y float32) {
	glyphRune := rest.GetSMUFLRune()
	glyphName := rest.GetSMUFLName()

	// Position rest in the middle of the staff
	restY := y - StaffSpacesToPixels(1) // Center on staff

	// Get bounding box if available
	bbox, hasBBox := sr.bboxMap[glyphName]

	// Draw the rest
	if hasBBox {
		drawGlyphWithBBox(sr.font, glyphRune, bbox, x, restY, 0)
	} else {
		rl.DrawTextEx(sr.font, string(glyphRune), rl.NewVector2(x, restY-(FontRenderSize/2)), FontRenderSize, 0.0, rl.Black)
	}
}

// RenderAccidental renders an accidental symbol
func (sr *ScoreRenderer) RenderAccidental(accidental string, x, y float32) {
	var glyphRune rune

	switch accidental {
	case "sharp":
		if r, err := sr.smuflMetadata.GetGlyphRune("accidentalSharp"); err == nil {
			glyphRune = r
		}
	case "flat":
		if r, err := sr.smuflMetadata.GetGlyphRune("accidentalFlat"); err == nil {
			glyphRune = r
		}
	case "natural":
		if r, err := sr.smuflMetadata.GetGlyphRune("accidentalNatural"); err == nil {
			glyphRune = r
		}
	default:
		return // Unknown accidental
	}

	if glyphRune != 0 {
		color := rl.Black
		if sr.elementCounter == sr.activeIndex {
			color = rl.Red
		}
		rl.DrawTextEx(sr.font, string(glyphRune), rl.NewVector2(x, y-(FontRenderSize/2)), FontRenderSize, 0.0, color)
	}
}

// BuildSchedule builds a linear playback schedule from the score
func (sr *ScoreRenderer) BuildSchedule(score *Score) {
	sr.schedule = []scheduledEvent{}
	current := 0.0
	secondsPerBeat := 60.0 / float64(score.Tempo)
	for _, m := range score.Measures {
		beats := m.ElementBeats()
		for i, e := range m.Elements {
			durSec := float64(beats[i]) * secondsPerBeat
			se := scheduledEvent{start: current, end: current + durSec}
			switch n := e.(type) {
			case Note:
				se.isNote = true
				se.pitch = n.Pitch
			}
			sr.schedule = append(sr.schedule, se)
			current += durSec
		}
	}
}

// StartPlayback starts playback from t=0
func (sr *ScoreRenderer) StartPlayback() {
	sr.startTime = rl.GetTime()
	sr.playing = true
	sr.activeIndex = -1
	sr.samplePhase = 0
}

func midiToFreq(m int) float64 {
	return 2 * 440.0 * math.Pow(2, (float64(m)-69.0)/12.0)
}

// updatePlayback updates activeIndex and feeds audio stream
func (sr *ScoreRenderer) updatePlayback() {
	if !sr.playing || len(sr.schedule) == 0 {
		return
	}
	t := rl.GetTime() - sr.startTime
	// find active index
	idx := -1
	for i, ev := range sr.schedule {
		if t >= ev.start && t < ev.end {
			idx = i
			break
		}
	}
	sr.activeIndex = idx
	// audio
	if !sr.audioReady {
		return
	}
	if rl.IsAudioStreamProcessed(sr.stream) {
		bufferSize := 512
		samples := make([]float32, bufferSize)
		freq := 0.0
		if idx >= 0 && sr.schedule[idx].isNote {
			freq = midiToFreq(sr.schedule[idx].pitch)
		}
		for i := 0; i < bufferSize; i++ {
			var sample float32 = 0
			if freq > 0 {
				sr.samplePhase += 2 * math.Pi * freq / float64(sr.sampleRate)
				sample = float32(0.2 * math.Sin(sr.samplePhase))
			}
			samples[i] = sample
		}
		rl.UpdateAudioStream(sr.stream, samples)
	}
}

// GetNoteStaffPosition calculates the staff line position for a MIDI note
// This is a simplified version - a full implementation would consider clefs
func GetNoteStaffPosition(midiNote int) int {
	// For treble clef, middle C (MIDI 60) is on a ledger line below the staff
	// Staff lines are: E4(64)=0, F4(65)=0.5, G4(67)=1, A4(69)=1.5, B4(71)=2, C5(72)=2.5, D5(74)=3, E5(76)=3.5, F5(77)=4

	// Simplified mapping for treble clef
	switch midiNote {
	case 60:
		return -2 // C4 - ledger line below staff
	case 62:
		return -1 // D4 - below staff
	case 64:
		return 0 // E4 - bottom staff line
	case 65:
		return 0 // F4 - between bottom and second line
	case 67:
		return 1 // G4 - second line
	case 69:
		return 2 // A4 - between second and third line
	case 71:
		return 3 // B4 - third line (middle)
	case 72:
		return 3 // C5 - between third and fourth line
	case 74:
		return 4 // D5 - fourth line
	case 76:
		return 5 // E5 - top staff line
	case 77:
		return 5 // F5 - between top line and above
	case 79:
		return 6 // G5 - above staff
	default:
		// Approximate mapping for other notes
		return (midiNote - 64) / 2 // E4 = staff line 0, each whole step up = +1 staff position
	}
}
