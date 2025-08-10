package music

// NoteValue represents note durations
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
	Number        int
	Elements      []MusicElement
	TimeSignature TimeSignature
}

// MusicElement interface implemented by Note and Rest
type MusicElement interface {
	GetDuration() NoteValue
}

// Note represents a musical note
type Note struct {
	Pitch      int // MIDI note number, e.g., 60 = middle C
	Duration   NoteValue
	StaffLine  int    // Position on staff, 0 = bottom line, etc.
	Accidental string // "", "sharp", "flat", "natural"
}

// Rest represents a musical rest
type Rest struct {
	Duration NoteValue
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

// AddNote adds a note to the measure
func (m *Measure) AddNote(pitch int, dur NoteValue, staffLine int, accidental string) {
	n := Note{
		Pitch:      pitch,
		Duration:   dur,
		StaffLine:  staffLine,
		Accidental: accidental,
	}
	m.Elements = append(m.Elements, n)
}

// AddRest adds a rest to the measure
func (m *Measure) AddRest(dur NoteValue) {
	r := Rest{Duration: dur}
	m.Elements = append(m.Elements, r)
}

// GetDuration returns duration of a Note
func (n Note) GetDuration() NoteValue {
	return n.Duration
}

// GetDuration returns duration of a Rest
func (r Rest) GetDuration() NoteValue {
	return r.Duration
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

func (m *Measure) ElementBeats() []float32 {
	beats := make([]float32, 0, len(m.Elements))
	for _, e := range m.Elements {
		q := durationQuarters(e.GetDuration())
		b := q * float32(m.TimeSignature.Denominator) / 4.0
		beats = append(beats, b)
	}
	return beats
}

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
// For example, if denominator=4 (quarter note gets the beat), an eighth note = 0.5 beats.
func noteValueToBeats(nv NoteValue, denominator int) float32 {
	// Convert NoteValue to quarter note units
	var quarterUnits float32
	switch nv {
	case WholeNote:
		quarterUnits = 4.0
	case HalfNote:
		quarterUnits = 2.0
	case QuarterNote:
		quarterUnits = 1.0
	case EighthNote:
		quarterUnits = 0.5
	case SixteenthNote:
		quarterUnits = 0.25
	case ThirtySecondNote:
		quarterUnits = 0.125
	case SixtyFourthNote:
		quarterUnits = 0.0625
	default:
		quarterUnits = 1.0
	}

	// Adjust based on denominator (e.g., denominator=4 means quarter note gets the beat)
	return quarterUnits * float32(denominator) / 4.0
}

func (n Note) GetSMUFLName() string {
	switch n.Duration {
	case WholeNote:
		return "noteWhole"
	case HalfNote:
		return "noteHalfUp"
	case QuarterNote:
		return "noteQuarterUp"
	case EighthNote:
		return "note8thUp"
	case SixteenthNote:
		return "note16thUp"
	case ThirtySecondNote:
		return "note32ndUp"
	case SixtyFourthNote:
		return "note64thUp"
	default:
		return "noteQuarterUp" // fallback
	}
}
