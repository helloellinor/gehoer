package main

import "strings"

// NoteSystem represents different note naming systems
type NoteSystem string

const (
	International NoteSystem = "international" // C D E F G A B
	Nynorsk      NoteSystem = "nynorsk"       // C D E F G A H
	German       NoteSystem = "german"        // C D E F G A H (with B for Bb)
)

// KeySignature represents a key signature with its tonic and mode
type KeySignature struct {
	Tonic string // Root note name (e.g., "C", "fiss", "b")
	Mode  string // "dur", "moll", "dorian", etc.
}

// String returns the Norwegian key signature name
func (k KeySignature) String() string {
	return k.Tonic + "-" + k.Mode
}

// NoteName represents a note with its letter name and accidental
type NoteName struct {
	Letter     string // C, D, E, F, G, A, H, B (B is used for Hb in Norwegian)
	Accidental string // "", "ess", "iss", "essess", "ississ"
}

// String returns the string representation of the note name
func (n NoteName) String() string {
	return n.Letter + n.Accidental
}

// Localization holds all localized strings and note systems
type Localization struct {
	Language     string
	NoteSystem   NoteSystem
	Terms        map[string]string
	NoteNames    map[int]string // MIDI note number to localized name
	KeySignature KeySignature   // Current key context
	KeySignatures []KeySignature // All available key signatures
}

// NewNynorskLocalization creates a Norwegian Nynorsk localization
func NewNynorskLocalization(keyTonic, keyMode string) *Localization {
	loc := &Localization{
		Language:   "nynorsk",
		NoteSystem: Nynorsk,
		KeySignature: KeySignature{
			Tonic: keyTonic,
			Mode:  keyMode,
		},
		Terms: map[string]string{
			// Basic musical terms
			"note":           "tone",
			"chord":          "akkord",
			"scale":          "skala",
			"key":            "toneart",
			"major":          "dur",
			"minor":          "moll",
			"tempo":          "tempo",
			"rhythm":         "rytme",
			"measure":        "takt",
			"staff":          "notesystem",
			"clef":           "nøkkel",
			"time_signature": "taktart",
			"key_signature":  "forteikn",
			
			// Modes (toneartar)
			"ionian":     "ionisk",      // same as major/dur
			"dorian":     "dorisk",
			"phrygian":   "frygisk",
			"lydian":     "lydisk",
			"mixolydian": "miksisk",     // or miksolydisk
			"aeolian":    "æolisk",      // same as minor/moll
			"locrian":    "lokrisk",
			
			// Note values
			"whole_note":         "heilnote",
			"half_note":          "halvnote", 
			"quarter_note":       "fjerdedelsnote",
			"eighth_note":        "åttedelsnote",
			"sixteenth_note":     "sekstendelsnote",
			"thirty_second_note": "trettitodelelsnote",
			"sixty_fourth_note":  "sekstifjerdelelsnote",
			
			// Accidentals
			"sharp":         "krysz",
			"flat":          "b",
			"natural":       "oppløysningsteikn",
			"double_sharp":  "dobbeltkrysz",
			"double_flat":   "dobbelt-b",
			"iss":           "iss", // Norwegian sharp suffix
			"ess":           "ess", // Norwegian flat suffix
			
			// Octave designations (Norwegian system)
			"subcontra":     "subkontra",
			"contra":        "kontra",
			"great":         "store", 
			"small":         "vesle",
			"one_line":      "einstroken", // c¹
			"two_line":      "tostroken",  // c²
			"three_line":    "trestroken", // c³
			"four_line":     "firestroken", // c⁴
			"five_line":     "femstroken", // c⁵
			
			// Intervals - basic
			"unison":        "prim",
			"second":        "sekund",
			"third":         "ters",
			"fourth":        "kvart",
			"fifth":         "kvint",
			"sixth":         "sekst",
			"seventh":       "septim",
			"octave":        "oktav",
			
			// Interval qualities (Norwegian)
			"diminished":         "forminska",
			"interval_minor":     "vesle",     
			"interval_major":     "store",     
			"perfect":            "rein",
			"augmented":          "forstørra",
			"doubly_diminished":  "dobbelforminska",
			"doubly_augmented":   "dobbelforstørra",
			
			// Specific intervals
			"minor_second":     "vesle sekund",
			"major_second":     "store sekund", 
			"minor_third":      "vesle ters",
			"major_third":      "store ters",
			"perfect_fourth":   "rein kvart",
			"tritone":          "tritonus",
			"diminished_fifth": "forminska kvint",
			"perfect_fifth":    "rein kvint",
			"minor_sixth":      "vesle sekst",
			"major_sixth":      "store sekst",
			"minor_seventh":    "vesle septim",
			"major_seventh":    "store septim",
			"perfect_octave":   "rein oktav",
			
			// Extended intervals
			"ninth":              "none",
			"eleventh":           "undecim",
			"thirteenth":         "terdecim",
			"minor_ninth":        "vesle none",
			"major_ninth":        "store none",
			"perfect_eleventh":   "rein undecim",
			"sharp_eleventh":     "forstørra undecim",
			"minor_thirteenth":   "vesle terdecim",
			"major_thirteenth":   "store terdecim",
		},
		NoteNames: make(map[int]string),
	}
	
	// Initialize all key signatures
	loc.initializeKeySignatures()
	
	return loc
}

// initializeKeySignatures creates all standard key signatures
func (l *Localization) initializeKeySignatures() {
	l.KeySignatures = []KeySignature{
		// Major keys (dur)
		{"C", "dur"},      // C major - no accidentals
		{"G", "dur"},      // G major - 1 sharp (fiss)
		{"D", "dur"},      // D major - 2 sharps (fiss, ciss)
		{"A", "dur"},      // A major - 3 sharps (fiss, ciss, giss)
		{"E", "dur"},      // E major - 4 sharps (fiss, ciss, giss, diss)
		{"H", "dur"},      // B major - 5 sharps (fiss, ciss, giss, diss, aiss)
		{"fiss", "dur"},   // F# major - 6 sharps
		{"ciss", "dur"},   // C# major - 7 sharps
		
		{"F", "dur"},      // F major - 1 flat (b)
		{"b", "dur"},      // Bb major - 2 flats (b, ess)
		{"ess", "dur"},    // Eb major - 3 flats (b, ess, ass)
		{"ass", "dur"},    // Ab major - 4 flats (b, ess, ass, dess)
		{"dess", "dur"},   // Db major - 5 flats (b, ess, ass, dess, gess)
		{"gess", "dur"},   // Gb major - 6 flats
		{"cess", "dur"},   // Cb major - 7 flats
		
		// Minor keys (moll)
		{"a", "moll"},     // A minor - no accidentals (relative to C major)
		{"e", "moll"},     // E minor - 1 sharp (fiss)
		{"h", "moll"},     // B minor - 2 sharps (fiss, ciss)
		{"fiss", "moll"},  // F# minor - 3 sharps (fiss, ciss, giss)
		{"ciss", "moll"},  // C# minor - 4 sharps (fiss, ciss, giss, diss)
		{"giss", "moll"},  // G# minor - 5 sharps (fiss, ciss, giss, diss, aiss)
		{"diss", "moll"},  // D# minor - 6 sharps
		{"aiss", "moll"},  // A# minor - 7 sharps
		
		{"d", "moll"},     // D minor - 1 flat (b)
		{"g", "moll"},     // G minor - 2 flats (b, ess)
		{"c", "moll"},     // C minor - 3 flats (b, ess, ass)
		{"f", "moll"},     // F minor - 4 flats (b, ess, ass, dess)
		{"b", "moll"},     // Bb minor - 5 flats (b, ess, ass, dess, gess)
		{"ess", "moll"},   // Eb minor - 6 flats
		{"ass", "moll"},   // Ab minor - 7 flats
		
		// Church modes on C (examples)
		{"C", "ionisk"},    // C Ionian (same as C major)
		{"C", "dorisk"},    // C Dorian
		{"C", "frygisk"},   // C Phrygian
		{"C", "lydisk"},    // C Lydian
		{"C", "miksisk"},   // C Mixolydian
		{"C", "æolisk"},    // C Aeolian (same as C minor)
		{"C", "lokrisk"},   // C Locrian
		
		// Common modes in other keys
		{"D", "dorisk"},    // D Dorian (popular in folk music)
		{"E", "frygisk"},   // E Phrygian
		{"F", "lydisk"},    // F Lydian
		{"G", "miksisk"},   // G Mixolydian
	}
}

// generateNynorskNoteNames creates a map from MIDI note numbers to Norwegian note names
func (l *Localization) generateNynorskNoteNames() {
	// Chromatic notes with correct Norwegian naming
	chromaticNotes := []struct {
		name     string
		semitone int
	}{
		{"C", 0},
		{"ciss", 1},    // C# 
		{"D", 2},
		{"diss", 3},    // D# 
		{"E", 4},
		{"F", 5},
		{"fiss", 6},    // F#
		{"G", 7},
		{"giss", 8},    // G#
		{"A", 9},
		{"aiss", 10},   // A# 
		{"H", 11},      // B in international
	}
	
	// Generate note names for MIDI range 0-127
	for midi := 0; midi <= 127; midi++ {
		octave := (midi / 12) - 1
		semitone := midi % 12
		
		// Get base note name, considering key signature
		noteName := l.getNoteNameInKey(chromaticNotes[semitone].name, semitone)
		
		// Add octave designation using Norwegian system
		octaveDesignation := l.getOctaveDesignation(octave, noteName)
		
		l.NoteNames[midi] = noteName + octaveDesignation
	}
}

// getNoteNameInKey returns the appropriate note name considering the key signature
func (l *Localization) getNoteNameInKey(defaultName string, semitone int) string {
	// For semitone 10 (A#/Bb), always prefer "b" in Norwegian
	if semitone == 10 {
		return "b"
	}
	
	// Consider key signature for enharmonic choices
	switch l.KeySignature.Tonic {
	case "C", "G", "D", "A", "E", "H", "fiss", "ciss": // Sharp keys
		// Prefer sharp names where applicable
		return defaultName
	case "F", "b", "ess", "ass", "dess", "gess", "cess": // Flat keys  
		// Convert to flat equivalents where applicable
		switch semitone {
		case 1:
			return "dess" // Db instead of ciss
		case 3:
			return "ess"  // Eb instead of diss
		case 6:
			return "gess" // Gb instead of fiss
		case 8:
			return "ass"  // Ab instead of giss
		case 10:
			return "b"    // Always "b" in Norwegian
		}
	}
	
	return defaultName
}

// getOctaveDesignation returns the Norwegian octave designation
func (l *Localization) getOctaveDesignation(octave int, noteName string) string {
	// Norwegian octave system based on C
	// C4 = einstroken C (c¹), etc.
	switch octave {
	case 0:
		return "" // subkontra (very rare)
	case 1:
		return "" // kontra 
	case 2:
		return "" // store
	case 3:
		return "" // vesle
	case 4:
		return "¹" // einstroken
	case 5:
		return "²" // tostroken
	case 6:
		return "³" // trestroken
	case 7:
		return "⁴" // firestroken
	case 8:
		return "⁵" // femstroken
	default:
		if octave > 8 {
			return "⁺" + string(rune(octave-8+'0'))
		}
		return ""
	}
}

// SetKeySignature updates the current key signature and regenerates note names
func (l *Localization) SetKeySignature(tonic, mode string) {
	l.KeySignature.Tonic = tonic
	l.KeySignature.Mode = mode
	l.generateNynorskNoteNames()
}

// GetNoteName returns the localized name for a MIDI note number
func (l *Localization) GetNoteName(midiNote int) string {
	if len(l.NoteNames) == 0 {
		l.generateNynorskNoteNames()
	}
	if name, exists := l.NoteNames[midiNote]; exists {
		return name
	}
	return "ukjent" // "unknown" in Nynorsk
}

// GetTerm returns a localized term
func (l *Localization) GetTerm(key string) string {
	if term, exists := l.Terms[key]; exists {
		return term
	}
	return key // Return the key if no translation exists
}

// ConvertNoteToMIDI converts a Norwegian note name string to MIDI number
func (l *Localization) ConvertNoteToMIDI(noteName string, octave int) int {
	note := strings.ToLower(strings.TrimSpace(noteName))
	
	var semitone int
	switch note {
	// Natural notes
	case "c":
		semitone = 0
	case "d":
		semitone = 2
	case "e":
		semitone = 4
	case "f":
		semitone = 5
	case "g":
		semitone = 7
	case "a":
		semitone = 9
	case "h":
		semitone = 11
	case "b": // This is Hb (H flat) in Norwegian!
		semitone = 10
		
	// Sharp variants (iss)
	case "ciss":
		semitone = 1
	case "diss":
		semitone = 3
	case "eiss":
		semitone = 5 // E# = F
	case "fiss":
		semitone = 6
	case "giss":
		semitone = 8
	case "aiss":
		semitone = 10
	case "hiss":
		semitone = 0 // H# = C
		
	// Flat variants (ess)
	case "cess":
		semitone = 11 // Cb = B/H
	case "dess":
		semitone = 1 // Db = C#
	case "eess":
		semitone = 3 // Eb = D#
	case "fess":
		semitone = 4 // Fb = E
	case "gess":
		semitone = 6 // Gb = F#
	case "ass":
		semitone = 8 // Ab = G#
	// Note: "b" is used instead of "hess" for Hb in Norwegian
		
	default:
		return -1 // Invalid note
	}
	
	return (octave+1)*12 + semitone
}

// GetIntervalName returns the Norwegian name for an interval in semitones
func (l *Localization) GetIntervalName(semitones int) string {
	switch semitones {
	case 0:
		return l.GetTerm("unison") // "prim"
	case 1:
		return l.GetTerm("minor_second") // "vesle sekund"
	case 2:
		return l.GetTerm("major_second") // "store sekund"
	case 3:
		return l.GetTerm("minor_third") // "vesle ters"
	case 4:
		return l.GetTerm("major_third") // "store ters"
	case 5:
		return l.GetTerm("perfect_fourth") // "rein kvart"
	case 6:
		return l.GetTerm("tritone") // "tritonus"
	case 7:
		return l.GetTerm("perfect_fifth") // "rein kvint"
	case 8:
		return l.GetTerm("minor_sixth") // "vesle sekst"
	case 9:
		return l.GetTerm("major_sixth") // "store sekst"
	case 10:
		return l.GetTerm("minor_seventh") // "vesle septim"
	case 11:
		return l.GetTerm("major_seventh") // "store septim"
	case 12:
		return l.GetTerm("perfect_octave") // "rein oktav"
	default:
		if semitones > 12 {
			// Compound intervals
			octaves := semitones / 12
			remainder := semitones % 12
			baseInterval := l.GetIntervalName(remainder)
			if octaves == 1 {
				return baseInterval + " + " + l.GetTerm("octave")
			}
			return baseInterval + " + " + string(rune(octaves+'0')) + " " + l.GetTerm("octave")
		}
		return "ukjent intervall"
	}
}

// GetKeySignatures returns all available key signatures
func (l *Localization) GetKeySignatures() []KeySignature {
	return l.KeySignatures
}

// Global localization instance
var Loc *Localization

// InitLocalization initializes the global localization with a key signature
func InitLocalization(keyTonic, keyMode string) {
	Loc = NewNynorskLocalization(keyTonic, keyMode)
}
