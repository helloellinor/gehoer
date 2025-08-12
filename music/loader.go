package music

import (
	"encoding/json"
	"os"
)

type JSONElement struct {
	Type       string `json:"type"`
	Pitch      int    `json:"pitch,omitempty"`
	Duration   string `json:"duration"`
	StaffLine  int    `json:"staff_line"`
	Accidental string `json:"accidental,omitempty"`
}

type JSONMeasure struct {
	Number   int           `json:"number"`
	Elements []JSONElement `json:"elements"`
}

type JSONStaff struct {
	Clef     string        `json:"clef"`
	Measures []JSONMeasure `json:"measures"`
}

type JSONScore struct {
	Title        string `json:"title"`
	Composer     string `json:"composer"`
	Instrument   string `json:"instrument,omitempty"` // "treble", "bass", "piano"
	KeySignature struct {
		Tonic string `json:"tonic"`
		Mode  string `json:"mode"`
	} `json:"key_signature"`
	TimeSignature struct {
		Numerator   int `json:"numerator"`
		Denominator int `json:"denominator"`
	} `json:"time_signature"`
	Tempo    int           `json:"tempo"`
	Staves   []JSONStaff   `json:"staves,omitempty"`   // Multi-staff format
	Measures []JSONMeasure `json:"measures,omitempty"` // Single-staff format (legacy)
}

func LoadScoreFromJSON(path string) (*Score, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var js JSONScore
	if err := json.Unmarshal(data, &js); err != nil {
		return nil, err
	}

	score := NewScore(js.Title, js.Composer, js.KeySignature.Tonic, js.KeySignature.Mode, js.TimeSignature.Numerator, js.TimeSignature.Denominator, js.Tempo)

	// Set instrument type (default to treble if not specified)
	if js.Instrument != "" {
		score.Instrument = js.Instrument
	} else {
		score.Instrument = InstrumentTreble
	}

	// Handle multi-staff format
	if len(js.Staves) > 0 {
		for _, jsonStaff := range js.Staves {
			staff := &Staff{
				Clef:     jsonStaff.Clef,
				Measures: make([]*Measure, 0),
			}

			for _, jm := range jsonStaff.Measures {
				measure := &Measure{
					Number:        jm.Number,
					Elements:      make([]MusicElement, 0),
					TimeSignature: score.TimeSignature,
				}

				for _, elem := range jm.Elements {
					switch elem.Type {
					case "note":
						dur := parseDuration(elem.Duration)
						measure.AddNote(&Note{
							Pitch:      elem.Pitch,
							Duration:   dur,
							StaffLine:  elem.StaffLine,
							Accidental: elem.Accidental,
						})
					case "rest":
						dur := parseDuration(elem.Duration)
						measure.AddRest(&Rest{
							Duration: dur,
						})
					}
				}
				staff.Measures = append(staff.Measures, measure)
			}
			score.Staves = append(score.Staves, staff)
		}
	} else {
		// Handle legacy single-staff format
		for _, jm := range js.Measures {
			measure := score.AddMeasure(nil) // Use default time signature or extend for per-measure
			for _, elem := range jm.Elements {
				switch elem.Type {
				case "note":
					dur := parseDuration(elem.Duration)
					measure.AddNote(&Note{
						Pitch:      elem.Pitch,
						Duration:   dur,
						StaffLine:  elem.StaffLine,
						Accidental: elem.Accidental,
					})
				case "rest":
					dur := parseDuration(elem.Duration)
					measure.AddRest(&Rest{
						Duration: dur,
					})
				}
			}
		}
	}

	return score, nil
}
