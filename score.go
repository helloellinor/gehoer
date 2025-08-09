package main

import (
	"fmt"
	"time"
)

// TimeSignature represents a time signature like 4/4, 3/4, etc.
type TimeSignature struct {
	Numerator   int // beats per measure
	Denominator int // note value that gets the beat (4 = quarter note)
}

// String returns the time signature as a string like "4/4"
func (ts TimeSignature) String() string {
	return fmt.Sprintf("%d/%d", ts.Numerator, ts.Denominator)
}

// NoteValue represents different note durations
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

// String returns the Norwegian name for the note value
func (nv NoteValue) String() string {
	if Loc == nil {
		return "ukjent"
	}
	
	switch nv {
	case WholeNote:
		return Loc.GetTerm("whole_note")
	case HalfNote:
		return Loc.GetTerm("half_note")
	case QuarterNote:
		return Loc.GetTerm("quarter_note")
	case EighthNote:
		return Loc.GetTerm("eighth_note")
	case SixteenthNote:
		return Loc.GetTerm("sixteenth_note")
	case ThirtySecondNote:
		return Loc.GetTerm("thirty_second_note")
	case SixtyFourthNote:
		return Loc.GetTerm("sixty_fourth_note")
	default:
		return "ukjent"
	}
}

// GetSMUFLRune returns the Unicode codepoint for the SMUFL glyph
func (nv NoteValue) GetSMUFLRune() rune {
	switch nv {
	case WholeNote:
		return '\uE0A2' // noteWhole
	case HalfNote:
		return '\uE0A3' // noteHalfUp
	case QuarterNote:
		return '\uE0A4' // noteQuarterUp
	default:
		return '\uE0A4' // Default to quarter note head
	}
}

// GetFlagRune returns the flag glyph for notes shorter than a quarter note
func (nv NoteValue) GetFlagRune() rune {
	switch nv {
	case EighthNote:
		return '\uE240' // flag8thUp
	case SixteenthNote:
		return '\uE242' // flag16thUp
	case ThirtySecondNote:
		return '\uE244' // flag32ndUp
	case SixtyFourthNote:
		return '\uE246' // flag64thUp
	default:
		return 0 // No flag for longer notes
	}
}

// GetSMUFLName returns the SMUFL glyph name for metadata lookup
func (nv NoteValue) GetSMUFLName() string {
	switch nv {
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
		return "noteQuarterUp"
	}
}

// Note represents a single musical note
type Note struct {
	Pitch     int       // MIDI note number
	Duration  NoteValue // Note duration
	StaffLine int       // Position on staff (0 = bottom line, negative = below staff)
	Accidental string   // "", "sharp", "flat", "natural"
}

// String returns the Norwegian note name
func (n Note) String() string {
	if Loc == nil {
		return fmt.Sprintf("Note(%d)", n.Pitch)
	}
	return Loc.GetNoteName(n.Pitch)
}

// Rest represents a musical rest
type Rest struct {
	Duration NoteValue // Rest duration
}

// GetSMUFLRune returns the SMUFL glyph for the rest
func (r Rest) GetSMUFLRune() rune {
	switch r.Duration {
	case WholeNote:
		return '\uE4E1' // restWhole
	case HalfNote:
		return '\uE4E2' // restHalf
	case QuarterNote:
		return '\uE4E5' // restQuarter
	case EighthNote:
		return '\uE4E6' // rest8th
	case SixteenthNote:
		return '\uE4E7' // rest16th
	case ThirtySecondNote:
		return '\uE4E8' // rest32nd
	case SixtyFourthNote:
		return '\uE4E9' // rest64th
	default:
		return '\uE4E5' // Default to quarter rest
	}
}

// GetSMUFLName returns the SMUFL glyph name for metadata lookup
func (r Rest) GetSMUFLName() string {
	switch r.Duration {
	case WholeNote:
		return "restWhole"
	case HalfNote:
		return "restHalf"
	case QuarterNote:
		return "restQuarter"
	case EighthNote:
		return "rest8th"
	case SixteenthNote:
		return "rest16th"
	case ThirtySecondNote:
		return "rest32nd"
	case SixtyFourthNote:
		return "rest64th"
	default:
		return "restQuarter"
	}
}

// MusicElement represents either a note or a rest
type MusicElement interface {
	GetDuration() NoteValue
}

func (n Note) GetDuration() NoteValue { return n.Duration }
func (r Rest) GetDuration() NoteValue { return r.Duration }

// Measure represents one measure of music
type Measure struct {
	Elements      []MusicElement // Notes and rests in this measure
	TimeSignature TimeSignature  // Time signature for this measure
	Number        int            // Measure number in the score
}

// durationQuarters returns duration in quarter-note units
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

// ElementBeats returns each element's duration measured in beats of the measure's denominator
func (m *Measure) ElementBeats() []float32 {
	beats := make([]float32, 0, len(m.Elements))
	for _, e := range m.Elements {
		q := durationQuarters(e.GetDuration())
		b := q * float32(m.TimeSignature.Denominator) / 4.0
		beats = append(beats, b)
	}
	return beats
}

// ElementPositions returns x positions for elements spaced proportionally by duration
// width: total width in pixels, margins: left and right margins in pixels
func (m *Measure) ElementPositions(width float32, leftMargin float32, rightMargin float32) []float32 {
	beats := m.ElementBeats()
	var totalBeats float32
	for _, b := range beats { totalBeats += b }
	usable := width - (leftMargin + rightMargin)
	positions := make([]float32, len(beats))
	acc := float32(0)
	for i := range beats {
		// place element at start position of its time slice
		positions[i] = leftMargin + (acc/totalBeats)*usable
		acc += beats[i]
	}
	return positions
}

// NewMeasure creates a new measure with the given time signature
func NewMeasure(number int, timeSignature TimeSignature) *Measure {
	return &Measure{
		Elements:      make([]MusicElement, 0),
		TimeSignature: timeSignature,
		Number:        number,
	}
}

// AddNote adds a note to the measure
func (m *Measure) AddNote(pitch int, duration NoteValue, staffLine int, accidental string) {
	note := Note{
		Pitch:      pitch,
		Duration:   duration,
		StaffLine:  staffLine,
		Accidental: accidental,
	}
	m.Elements = append(m.Elements, note)
}

// AddRest adds a rest to the measure
func (m *Measure) AddRest(duration NoteValue) {
	rest := Rest{Duration: duration}
	m.Elements = append(m.Elements, rest)
}

// IsFull checks if the measure is full based on its time signature
func (m *Measure) IsFull() bool {
	totalDuration := 0
	denominator := m.TimeSignature.Denominator
	
	for _, element := range m.Elements {
		// Convert note value to duration units based on denominator
		duration := element.GetDuration()
		switch duration {
		case WholeNote:
			totalDuration += 4 * (4 / denominator)
		case HalfNote:
			totalDuration += 2 * (4 / denominator)
		case QuarterNote:
			totalDuration += 1 * (4 / denominator)
		case EighthNote:
			totalDuration += (4 / denominator) / 2
		case SixteenthNote:
			totalDuration += (4 / denominator) / 4
		case ThirtySecondNote:
			totalDuration += (4 / denominator) / 8
		case SixtyFourthNote:
			totalDuration += (4 / denominator) / 16
		}
	}
	
	return totalDuration >= m.TimeSignature.Numerator
}

// String returns a string representation of the measure
func (m *Measure) String() string {
	return fmt.Sprintf("Measure %d (%s): %d elements", 
		m.Number, m.TimeSignature.String(), len(m.Elements))
}

// Score represents a complete musical score
type Score struct {
	Title         string         // Title of the piece
	Composer      string         // Composer name
	KeySignature  KeySignature   // Key signature
	TimeSignature TimeSignature  // Default time signature
	Tempo         int            // BPM (beats per minute)
	Measures      []*Measure     // All measures in the score
	CreatedAt     time.Time      // When the score was created
}

// NewScore creates a new empty score
func NewScore(title, composer string, keyTonic, keyMode string, timeNum, timeDen, tempo int) *Score {
	return &Score{
		Title:    title,
		Composer: composer,
		KeySignature: KeySignature{
			Tonic: keyTonic,
			Mode:  keyMode,
		},
		TimeSignature: TimeSignature{
			Numerator:   timeNum,
			Denominator: timeDen,
		},
		Tempo:     tempo,
		Measures:  make([]*Measure, 0),
		CreatedAt: time.Now(),
	}
}

// AddMeasure adds a new measure to the score
func (s *Score) AddMeasure(timeSignature *TimeSignature) *Measure {
	measureNumber := len(s.Measures) + 1
	
	// Use provided time signature or default
	ts := s.TimeSignature
	if timeSignature != nil {
		ts = *timeSignature
	}
	
	measure := NewMeasure(measureNumber, ts)
	s.Measures = append(s.Measures, measure)
	return measure
}

// GetMeasure returns the measure at the given index (1-based)
func (s *Score) GetMeasure(number int) *Measure {
	if number < 1 || number > len(s.Measures) {
		return nil
	}
	return s.Measures[number-1]
}

// GetMeasureCount returns the total number of measures
func (s *Score) GetMeasureCount() int {
	return len(s.Measures)
}

// String returns a string representation of the score
func (s *Score) String() string {
	return fmt.Sprintf("\"%s\" by %s (%s, %s, %d BPM) - %d measures", 
		s.Title, s.Composer, s.KeySignature.String(), 
		s.TimeSignature.String(), s.Tempo, len(s.Measures))
}

// SetKeySignature updates the key signature and localization
func (s *Score) SetKeySignature(tonic, mode string) {
	s.KeySignature.Tonic = tonic
	s.KeySignature.Mode = mode
	
	// Update global localization if it exists
	if Loc != nil {
		Loc.SetKeySignature(tonic, mode)
	}
}
