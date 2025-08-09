package main

// units.go - SMUFL-compliant unit system
// Based on: https://w3c.github.io/smufl/latest/specification/scoring-metrics-glyph-registration.html
//
// SMUFL Specification Guidelines:
// - "Dividing the em in four provides an analogue for a five-line staff"
// - "one staff space = 0.25 em"
// - All measurements should be in terms of ems and staff spaces

const (
	// PRIMARY CONSTANT: Font size in pixels (our only visual choice)
	// This represents the em size - the fundamental unit in SMUFL
	EmSizePx = 100 // Adjust this to make everything bigger/smaller

	// Technical constants
	FontLoadMultiplier       = 16 // Load at high resolution to avoid pixelation
	GridSpacingInStaffSpaces = 2  // Grid lines every N staff spaces

	// Empirical correction factor (due to Raylib font scaling behavior)
	RaylibFontScaleFactor = 4.0 // Adjust if font doesn't match expected size
)

// SMUFL-compliant derived units
var (
	// Em size in pixels (SMUFL fundamental unit)
	EmPx = float32(EmSizePx)

	// Staff space size (SMUFL: exactly 0.25 em)
	StaffSpacePx = EmPx * 0.25

	// Font rendering size (adjusted for Raylib behavior)
	FontRenderSize = EmPx * RaylibFontScaleFactor

	// Font load size for HiDPI
	FontLoadSize = EmSizePx * FontLoadMultiplier

	// Grid spacing (aligned to staff spaces)
	GridSpacingPx = int(StaffSpacePx * GridSpacingInStaffSpaces)

	// Grid font size (readable relative to em size)
	GridFontSize = int32(EmSizePx / 15)

	// SMUFL bounding box scale (1 SMUFL unit = 1 staff space)
	SMUFLBBoxScale = StaffSpacePx
)

// SMUFL Unit Conversion Functions

// EmsToPixels converts em units to pixels
func EmsToPixels(ems float32) float32 {
	return ems * EmPx
}

// StaffSpacesToPixels converts staff space units to pixels
func StaffSpacesToPixels(staffSpaces float32) float32 {
	return staffSpaces * StaffSpacePx
}

// PixelsToEms converts pixels to em units
func PixelsToEms(pixels float32) float32 {
	return pixels / EmPx
}

// PixelsToStaffSpaces converts pixels to staff space units
func PixelsToStaffSpaces(pixels float32) float32 {
	return pixels / StaffSpacePx
}

// EmsToStaffSpaces converts em units to staff space units (should always be 4:1)
func EmsToStaffSpaces(ems float32) float32 {
	return ems * 4.0
}

// StaffSpacesToEms converts staff space units to em units (should always be 1:4)
func StaffSpacesToEms(staffSpaces float32) float32 {
	return staffSpaces * 0.25
}
