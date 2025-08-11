package units

import "gehoer/settings"

// Derived unit variables initialized on package load
var (
	FontLoadSizePx   = settings.MusicFontSizePx * settings.MusicFontLoadMultiplier
	FontRenderSizePx = EmPx * settings.RaylibFontScaleFactor

	EmPx         = float32(settings.MusicFontSizePx)
	StaffSpacePx = EmPx * 0.25

	GridSpacingPx  = int(StaffSpacePx * float32(settings.GridSpacingInStaffSpaces))
	GridFontSizePx = int32(settings.MusicFontSizePx / 15)
)

func StaffSpacesToPixels(staffSpaces float32) float32 {
	return staffSpaces * StaffSpacePx
}

func PixelsToStaffSpaces(pixels float32) float32 {
	return pixels / StaffSpacePx
}

func EmsToPixels(ems float32) float32 {
	return ems * EmPx
}

func PixelsToEms(pixels float32) float32 {
	return pixels / EmPx
}

func EmsToStaffSpaces(ems float32) float32 {
	return ems * 4.0
}

func StaffSpacesToEms(staffSpaces float32) float32 {
	return staffSpaces * 0.25
}
