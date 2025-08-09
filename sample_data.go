package main

// CreateLisaGikkTilSkolenScore creates the score for the Norwegian children's song
// "Lisa gikk til skolen" in C major, 4/4 time
func CreateLisaGikkTilSkolenScore() *Score {
	// Create the score in C major (C-dur), 4/4 time, moderate tempo
	score := NewScore("Lisa gikk til skolen", "Norsk barnesang", "C", "dur", 4, 4, 120)
	
	// The melody for "Lisa gikk til skolen":
	// Li-sa gikk til sko-len, sko-len, sko-len
	// Li-sa gikk til sko-len så fint i dag
	// Ut og hjem og ut og hjem og ut og hjem igjen
	// Li-sa gikk til sko-len så fint i dag
	
	// Note: Using MIDI note numbers where C4 = 60
	// C4=60, D4=62, E4=64, F4=65, G4=67, A4=69, B4=71, C5=72
	
	// Measure 1: "Li-sa gikk til" (C D E F)
	measure1 := score.AddMeasure(nil)
	measure1.AddNote(60, QuarterNote, GetNoteStaffPosition(60), "") // C4 - Li
	measure1.AddNote(62, QuarterNote, GetNoteStaffPosition(62), "") // D4 - sa  
	measure1.AddNote(64, QuarterNote, GetNoteStaffPosition(64), "") // E4 - gikk
	measure1.AddNote(65, QuarterNote, GetNoteStaffPosition(65), "") // F4 - til
	
	// Measure 2: "sko-len, sko-len" (G G F F)
	measure2 := score.AddMeasure(nil)
	measure2.AddNote(67, QuarterNote, GetNoteStaffPosition(67), "") // G4 - sko
	measure2.AddNote(67, QuarterNote, GetNoteStaffPosition(67), "") // G4 - len
	measure2.AddNote(65, QuarterNote, GetNoteStaffPosition(65), "") // F4 - sko  
	measure2.AddNote(65, QuarterNote, GetNoteStaffPosition(65), "") // F4 - len
	
	// Measure 3: "sko-len" (E D)
	measure3 := score.AddMeasure(nil)
	measure3.AddNote(64, HalfNote, 2, "")    // E4 - sko
	measure3.AddNote(62, HalfNote, 1, "")    // D4 - len
	
	// Measure 4: "Li-sa gikk til" (C D E F) 
	measure4 := score.AddMeasure(nil)
	measure4.AddNote(60, QuarterNote, 0, "") // C4 - Li
	measure4.AddNote(62, QuarterNote, 1, "") // D4 - sa
	measure4.AddNote(64, QuarterNote, 2, "") // E4 - gikk  
	measure4.AddNote(65, QuarterNote, 3, "") // F4 - til
	
	// Measure 5: "sko-len så" (G G A)
	measure5 := score.AddMeasure(nil)
	measure5.AddNote(67, QuarterNote, 4, "") // G4 - sko
	measure5.AddNote(67, QuarterNote, 4, "") // G4 - len
	measure5.AddNote(69, HalfNote, 5, "")    // A4 - så
	
	// Measure 6: "fint i dag" (G F E)
	measure6 := score.AddMeasure(nil)
	measure6.AddNote(67, QuarterNote, 4, "") // G4 - fint
	measure6.AddNote(65, QuarterNote, 3, "") // F4 - i
	measure6.AddNote(64, HalfNote, 2, "")    // E4 - dag
	
	// Measure 7: "Ut og hjem og" (E F G A)
	measure7 := score.AddMeasure(nil)
	measure7.AddNote(64, QuarterNote, 2, "") // E4 - Ut
	measure7.AddNote(65, QuarterNote, 3, "") // F4 - og
	measure7.AddNote(67, QuarterNote, 4, "") // G4 - hjem  
	measure7.AddNote(69, QuarterNote, 5, "") // A4 - og
	
	// Measure 8: "ut og hjem og" (A G F E)
	measure8 := score.AddMeasure(nil)
	measure8.AddNote(69, QuarterNote, 5, "") // A4 - ut
	measure8.AddNote(67, QuarterNote, 4, "") // G4 - og
	measure8.AddNote(65, QuarterNote, 3, "") // F4 - hjem
	measure8.AddNote(64, QuarterNote, 2, "") // E4 - og
	
	// Measure 9: "ut og hjem i-" (F G A A)  
	measure9 := score.AddMeasure(nil)
	measure9.AddNote(65, QuarterNote, 3, "") // F4 - ut
	measure9.AddNote(67, QuarterNote, 4, "") // G4 - og  
	measure9.AddNote(69, QuarterNote, 5, "") // A4 - hjem
	measure9.AddNote(69, QuarterNote, 5, "") // A4 - i
	
	// Measure 10: "-gjen" (G F)
	measure10 := score.AddMeasure(nil)
	measure10.AddNote(67, HalfNote, 4, "")   // G4 - gjen
	measure10.AddRest(HalfNote)              // Rest
	
	// Measure 11: "Li-sa gikk til" (C D E F) 
	measure11 := score.AddMeasure(nil)
	measure11.AddNote(60, QuarterNote, 0, "") // C4 - Li
	measure11.AddNote(62, QuarterNote, 1, "") // D4 - sa
	measure11.AddNote(64, QuarterNote, 2, "") // E4 - gikk
	measure11.AddNote(65, QuarterNote, 3, "") // F4 - til
	
	// Measure 12: "sko-len så" (G G A)
	measure12 := score.AddMeasure(nil)
	measure12.AddNote(67, QuarterNote, 4, "") // G4 - sko
	measure12.AddNote(67, QuarterNote, 4, "") // G4 - len
	measure12.AddNote(69, HalfNote, 5, "")    // A4 - så
	
	// Measure 13: "fint i" (G F)  
	measure13 := score.AddMeasure(nil)
	measure13.AddNote(67, HalfNote, 4, "")   // G4 - fint
	measure13.AddNote(65, HalfNote, 3, "")   // F4 - i
	
	// Measure 14: "dag" (C - whole note ending)
	measure14 := score.AddMeasure(nil)
	measure14.AddNote(60, WholeNote, 0, "")  // C4 - dag
	
	return score
}

// CreateSimpleCMajorScale creates a simple C major scale for testing
func CreateSimpleCMajorScale() *Score {
	score := NewScore("C-dur skala", "Øvelse", "C", "dur", 4, 4, 100)
	
	// C major scale: C D E F G A B C
	measure1 := score.AddMeasure(nil)
	measure1.AddNote(60, QuarterNote, 0, "") // C4
	measure1.AddNote(62, QuarterNote, 1, "") // D4
	measure1.AddNote(64, QuarterNote, 2, "") // E4
	measure1.AddNote(65, QuarterNote, 3, "") // F4
	
	measure2 := score.AddMeasure(nil)
	measure2.AddNote(67, QuarterNote, 4, "") // G4
	measure2.AddNote(69, QuarterNote, 5, "") // A4
	measure2.AddNote(71, QuarterNote, 6, "") // B4 (H in Norwegian)
	measure2.AddNote(72, QuarterNote, 7, "") // C5
	
	return score
}

// CreateNorwegianFolkSong creates a simple Norwegian folk song pattern
func CreateNorwegianFolkSong() *Score {
	// Simple folk song in G major
	score := NewScore("Norsk folkevise", "Tradisjonell", "G", "dur", 3, 4, 90)
	
	// Simple waltz-like pattern in 3/4 time
	// G major scale notes: G=67, A=69, B=71, C=72, D=74, E=76, F#=78, G=79
	
	measure1 := score.AddMeasure(nil)
	measure1.AddNote(67, QuarterNote, 0, "") // G4
	measure1.AddNote(69, QuarterNote, 1, "") // A4  
	measure1.AddNote(71, QuarterNote, 2, "") // H4
	
	measure2 := score.AddMeasure(nil)
	measure2.AddNote(72, HalfNote, 3, "")    // C5
	measure2.AddNote(71, QuarterNote, 2, "") // H4
	
	measure3 := score.AddMeasure(nil)
	measure3.AddNote(69, QuarterNote, 1, "") // A4
	measure3.AddNote(67, HalfNote, 0, "")    // G4
	
	measure4 := score.AddMeasure(nil)
	measure4.AddRest(HalfNote)               // Rest
	measure4.AddNote(74, QuarterNote, 4, "") // D5
	
	return score
}

// GetAllSampleScores returns all available sample scores
func GetAllSampleScores() map[string]*Score {
	return map[string]*Score{
		"lisa_gikk_til_skolen": CreateLisaGikkTilSkolenScore(),
		"c_major_scale":        CreateSimpleCMajorScale(),
		"norwegian_folk":       CreateNorwegianFolkSong(),
	}
}
