package main

import (
	"log"

	"gehoer/camera"
	"gehoer/fontloader"
	"gehoer/grid"
	"gehoer/metadata"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

func drawGlyphWithBBox(
	font rl.Font,
	r rune,
	bBox metadata.GlyphBBox,
	originX, originY float32,
	noteValueInStaffSpaces float32,
) {
	// Convert SMUFL bounding box coordinates (in staff spaces) to pixels
	bBoxWidth := StaffSpacesToPixels(float32(bBox.NE[0]) - float32(bBox.SW[0]))
	bBoxHeight := StaffSpacesToPixels(float32(bBox.NE[1]) - float32(bBox.SW[1]))
	staffSpaceOffset := StaffSpacesToPixels(noteValueInStaffSpaces)

	// Calculate bounding box position (SMUFL coordinates relative to baseline)
	rectX := originX + StaffSpacesToPixels(float32(bBox.SW[0]))
	rectY := originY - StaffSpacesToPixels(float32(bBox.NE[1])) - staffSpaceOffset

	rl.DrawRectangleLinesEx(rl.NewRectangle(rectX, rectY, bBoxWidth, bBoxHeight), 1.0, rl.Red)

	// Draw the glyph (positioned at baseline)
	glyphX := originX
	glyphY := originY - staffSpaceOffset - (FontRenderSize / 2)
	rl.DrawTextEx(font, string(r), rl.NewVector2(glyphX, glyphY), FontRenderSize, 0.0, rl.Black)
}

func main() {
	// Initialize Norwegian localization in C major
	InitLocalization("C", "dur")

	// DEBUG: Print our SMUFL unit system
	log.Printf("\n=== SMUFL UNIT SYSTEM ===")
	log.Printf("EmSizePx: %.1f px (1 em)", float32(EmSizePx))
	log.Printf("FontRenderSize: %.1f px (equals EmSizePx)", FontRenderSize)
	log.Printf("StaffSpacePx: %.1f px (0.25 em)", StaffSpacePx)
	log.Printf("SMUFLBBoxScale: %.1f px/staff-space (bounding box scaling)", SMUFLBBoxScale)
	log.Printf("FontLoadSize: %d px (font loading resolution)", FontLoadSize)
	log.Printf("GridSpacingPx: %d px", GridSpacingPx)
	log.Printf("GridFontSize: %d px", GridFontSize)

	// Print localization info
	log.Printf("\n=== LOKALISERING ===")
	log.Printf("Språk: %s", Loc.Language)
	log.Printf("Notesystem: %s", Loc.NoteSystem)
	log.Printf("Toneart: %s", Loc.KeySignature.String())
	log.Printf("C4 heter: %s", Loc.GetNoteName(60))
	log.Printf("A#4/Bb4 heter: %s", Loc.GetNoteName(70))
	log.Printf("Vesle sekund: %s", Loc.GetIntervalName(1))

	// Create sample scores
	sampleScores := GetAllSampleScores()
	lisaScore := sampleScores["lisa_gikk_til_skolen"]

	log.Printf("\n=== PARTITUR ===")
	log.Printf("Laga partitur: %s", lisaScore.String())
	for i, measure := range lisaScore.Measures[:3] { // Show first 3 measures
		log.Printf("Takt %d: %s", i+1, measure.String())
	}

	// Load SMuFL metadata from the official repository
	smuflMetadata, err := metadata.LoadSMuFLMetadata("external/smufl")
	if err != nil {
		log.Fatalf("Failed to load SMuFL metadata: %v", err)
	}

	// Print SMuFL info
	log.Printf("\n=== SMuFL METADATA ===")
	log.Printf("Loaded %d glyphs from official repository", len(smuflMetadata.Glyphs))
	log.Printf("Available ranges: %d", len(smuflMetadata.Ranges))
	log.Printf("Available classes: %d", len(smuflMetadata.Classes))

	// Load only the glyphs we need now: basic noteheads, clefs, accidentals, time signature digits, and flags used
	// Basic noteheads we render
	noteheadNames := []string{"noteheadWhole", "noteheadHalf", "noteheadBlack"}
	// Clefs and accidentals from metadata helpers
	clefs, _ := smuflMetadata.GetClefGlyphs()
	accidentals, _ := smuflMetadata.GetAccidentalGlyphs()
	// Time signature digits for Lisa gikk til skolen (4/4)
	timeSigNames := []string{"timeSig4"}
	// Basic flags for up/down (8th and 16th for now; extend if needed)
	flagNames := []string{"flag8thUp", "flag8thDown", "flag16thUp", "flag16thDown"}

	// Collect runes
	allGlyphRunes := []rune{}
	for _, name := range noteheadNames {
		if r, err := smuflMetadata.GetGlyphRune(name); err == nil {
			allGlyphRunes = append(allGlyphRunes, r)
		} else {
			log.Printf("Warning: missing glyph %s: %v", name, err)
		}
	}
	for _, r := range clefs {
		allGlyphRunes = append(allGlyphRunes, r)
	}
	for _, r := range accidentals {
		allGlyphRunes = append(allGlyphRunes, r)
	}
	for _, name := range timeSigNames {
		if r, err := smuflMetadata.GetGlyphRune(name); err == nil {
			allGlyphRunes = append(allGlyphRunes, r)
		} else {
			log.Printf("Warning: missing glyph %s: %v", name, err)
		}
	}
	for _, name := range flagNames {
		if r, err := smuflMetadata.GetGlyphRune(name); err == nil {
			allGlyphRunes = append(allGlyphRunes, r)
		} else {
			log.Printf("Warning: missing glyph %s: %v", name, err)
		}
	}

	log.Printf("Loading %d selected glyphs (noteheads, clefs, accidentals, time signatures, flags)", len(allGlyphRunes))

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Gehør - Lisa gikk til skolen")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	font, err := fontloader.LoadSmuflFont("assets/fonts/Leland/Leland.otf", int32(FontLoadSize), allGlyphRunes)
	if err != nil {
		log.Fatalf("Font load error: %v", err)
	}
	defer rl.UnloadFont(font)

	// Create bounding box map from SMuFL metadata
	bboxMap := smuflMetadata.CreateGlyphBoundingBoxMap()
	// Load font-specific anchors and engraving defaults (generic loaders)
	fontMetaPath := "assets/fonts/Leland/leland_metadata.json"
	anchors, err := metadata.LoadFontAnchors(fontMetaPath)
	if err != nil {
		log.Printf("Warning: could not load font anchors: %v", err)
	}
	defaults, err := metadata.LoadFontEngravingDefaults(fontMetaPath)
	if err != nil {
		log.Printf("Warning: could not load engraving defaults: %v", err)
	}
	// Create score renderer
	scoreRenderer := NewScoreRenderer(smuflMetadata, font, bboxMap, anchors, defaults)
	// Expose renderer to staff helpers for thickness
	currentRenderer = scoreRenderer

	myGrid := grid.New(GridSpacingPx, 4000, 4000, GridFontSize)
	camCtrl := camera.NewCameraController(screenWidth, screenHeight)

	for !rl.WindowShouldClose() {
		// Update camera input & state
		camCtrl.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode2D(camCtrl.Camera)

		// Draw grid with labels
		myGrid.Draw(camCtrl.Camera)

		// Render the "Lisa gikk til skolen" score
		scoreRenderer.RenderScore(lisaScore, 50, 200)

		rl.EndMode2D()

		// UI text fixed on screen
		rl.DrawText("Arrow keys to pan, mouse drag to pan, scroll to zoom", 10, 10, 20, rl.DarkGray)
		rl.DrawText("Norwegian Children's Song: Lisa gikk til skolen", 10, 40, 16, rl.DarkBlue)

		rl.EndDrawing()
	}
}
