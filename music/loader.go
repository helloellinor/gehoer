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

type JSONScore struct {
	Title        string `json:"title"`
	Composer     string `json:"composer"`
	KeySignature struct {
		Tonic string `json:"tonic"`
		Mode  string `json:"mode"`
	} `json:"key_signature"`
	TimeSignature struct {
		Numerator   int `json:"numerator"`
		Denominator int `json:"denominator"`
	} `json:"time_signature"`
	Tempo    int           `json:"tempo"`
	Measures []JSONMeasure `json:"measures"`
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

	for _, jm := range js.Measures {
		measure := score.AddMeasure(nil) // Use default time signature or extend for per-measure
		for _, elem := range jm.Elements {
			switch elem.Type {
			case "note":
				dur := parseDuration(elem.Duration)
				measure.AddNote(elem.Pitch, dur, elem.StaffLine, elem.Accidental)
			case "rest":
				dur := parseDuration(elem.Duration)
				measure.AddRest(dur)
			}
		}
	}

	return score, nil
}

func parseDuration(dur string) NoteValue {
	switch dur {
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
	// add more if you like
	default:
		return QuarterNote // fallback
	}
}
