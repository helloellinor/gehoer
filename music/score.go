package music

// NoteValue represents note durations as an enum
type NoteValue int

const (
	WholeNote NoteValue = iota
	HalfNote
	QuarterNote
	EighthNote
	SixteenthNote
	ThirtySecondNote
	SixtyFourthNote
)

// Score represents a full musical score
type Score struct {
	Title         string
	Composer      string
	KeySignature  KeySignature
	TimeSignature TimeSignature
	Tempo         int
	Measures      []*Measure
}

// KeySignature stores tonic and mode info
type KeySignature struct {
	Tonic string
	Mode  string
}

// TimeSignature stores numerator and denominator
type TimeSignature struct {
	Numerator   int
	Denominator int
}

// Measure is one bar of music containing notes and rests
type Measure struct {
	Clef          string
	Number        int
	Elements      []MusicElement
	TimeSignature TimeSignature
}

// MusicElement interface implemented by Note and Rest
type MusicElement interface {
	GetDuration() NoteValue
	GlyphName() string
}

// Note represents a musical note
type Note struct {
	Pitch      int // MIDI note number, e.g., 60 = middle C
	Duration   NoteValue
	StaffLine  int    // Position on staff, 0 = bottom line, etc.
	Accidental string // "", "sharp", "flat", "natural"
}

func (n *Note) GetDuration() NoteValue {
	return n.Duration
}

func (n *Note) HasStem() bool {
	switch n.Duration {
	case WholeNote:
		return false
	default:
		return true
	}
}

func (n *Note) HasFlag() bool {
	switch n.Duration {
	case EighthNote, SixteenthNote, ThirtySecondNote, SixtyFourthNote:
		return true
	default:
		return false
	}
}

// NoteheadGlyphName returns the SMuFL glyph name for the notehead based on Duration
func (n *Note) NoteheadGlyphName() string {
	switch n.Duration {
	case WholeNote:
		return "noteheadWhole"
	case HalfNote:
		return "noteheadHalf"
	case QuarterNote, EighthNote:
		return "noteheadBlack"
	// Add more durations if needed
	default:
		return "noteheadBlack" // fallback
	}
}

// Rest represents a musical rest
type Rest struct {
	Duration NoteValue
}

func (r *Rest) GetDuration() NoteValue {
	return r.Duration
}

// NewScore creates a new score with initial values and empty measure slice
func NewScore(title, composer, tonic, mode string, timeNum, timeDen, tempo int) *Score {
	return &Score{
		Title:    title,
		Composer: composer,
		KeySignature: KeySignature{
			Tonic: tonic,
			Mode:  mode,
		},
		TimeSignature: TimeSignature{
			Numerator:   timeNum,
			Denominator: timeDen,
		},
		Tempo:    tempo,
		Measures: make([]*Measure, 0),
	}
}

// AddMeasure appends a new measure and returns it
func (s *Score) AddMeasure(ts *TimeSignature) *Measure {
	t := s.TimeSignature
	if ts != nil {
		t = *ts
	}
	measure := &Measure{
		Elements:      make([]MusicElement, 0),
		TimeSignature: t,
	}
	s.Measures = append(s.Measures, measure)
	return measure
}

func (m *Measure) AddNote(note *Note) {
	m.Elements = append(m.Elements, note)
}

func (m *Measure) AddRest(rest *Rest) {
	m.Elements = append(m.Elements, rest)
}

// ElementPositions returns x positions for elements spaced proportionally by duration
func (m *Measure) ElementPositions(width float32, leftMargin float32, rightMargin float32) []float32 {
	beats := m.ElementBeats()
	var totalBeats float32
	for _, b := range beats {
		totalBeats += b
	}
	usable := width - (leftMargin + rightMargin)
	positions := make([]float32, len(beats))
	acc := float32(0)
	for i := range beats {
		positions[i] = leftMargin + (acc/totalBeats)*usable
		acc += beats[i]
	}
	return positions
}

// ElementBeats returns beat lengths for each element in the measure
func (m *Measure) ElementBeats() []float32 {
	beats := make([]float32, 0, len(m.Elements))
	for _, e := range m.Elements {
		q := durationQuarters(e.GetDuration())
		b := q * float32(m.TimeSignature.Denominator) / 4.0
		beats = append(beats, b)
	}
	return beats
}

// durationQuarters converts NoteValue to quarter note units
func durationQuarters(nv NoteValue) float32 {
	switch nv {
	case WholeNote:
		return 4.0
	case HalfNote:
		return 2.0
	case QuarterNote:
		return 1.0
	case EighthNote:
		return 0.5
	case SixteenthNote:
		return 0.25
	case ThirtySecondNote:
		return 0.125
	case SixtyFourthNote:
		return 0.0625
	default:
		return 1.0
	}
}

// noteValueToBeats converts a NoteValue to beats relative to the measure's denominator.
func noteValueToBeats(nv NoteValue, denominator int) float32 {
	quarterUnits := durationQuarters(nv)
	return quarterUnits * float32(denominator) / 4.0
}

// parseDuration converts string duration to NoteValue enum
func parseDuration(s string) NoteValue {
	switch s {
	case "whole":
		return WholeNote
	case "half":
		return HalfNote
	case "quarter":
		return QuarterNote
	case "eighth":
		return EighthNote
	case "sixteenth":
		return SixteenthNote
	case "thirtysecond":
		return ThirtySecondNote
	case "sixtyfourth":
		return SixtyFourthNote
	default:
		return QuarterNote
	}
}

func (n *Note) GlyphName() string {
	return n.NoteheadGlyphName()
}

func (r *Rest) GlyphName() string {
	switch r.Duration {
	case WholeNote:
		return "restWhole"
	case HalfNote:
		return "restHalf"
	case QuarterNote:
		return "restQuarter"
	// add other durations as needed
	default:
		return "restQuarter"
	}
}
