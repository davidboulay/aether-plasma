package template

import "aether/internal/color"

// mixHex blends a toward b by weight (0 = a, 1 = b).
func mixHex(a, b string, weight float64) string {
	ca, cb := color.HexToRGB(a), color.HexToRGB(b)
	return color.RGBToHex(
		ca.R+(cb.R-ca.R)*weight,
		ca.G+(cb.G-ca.G)*weight,
		ca.B+(cb.B-ca.B)*weight,
	)
}

// readableOn returns black or white, whichever contrasts more with bg.
func readableOn(bg string) string {
	if color.ContrastRatio(bg, "#000000") >= color.ContrastRatio(bg, "#ffffff") {
		return "#000000"
	}
	return "#ffffff"
}

// plasmaVariableKeys lists every variable addPlasmaVariables may set.
var plasmaVariableKeys = []string{
	"kde_view", "kde_view_alt", "kde_window", "kde_window_alt",
	"kde_button", "kde_button_alt", "kde_header", "kde_header_inactive",
	"kde_tooltip", "kde_fg_inactive", "kde_selection", "kde_selection_alt",
	"kde_selection_fg", "kde_selection_fg_inactive", "kde_hover",
}

// addPlasmaVariables derives the surface shades a KDE color scheme needs.
// Shades step from the background toward the foreground, so they work the
// same way in light and dark mode (the Omarchy-era dark_bg/lighter_bg shades
// only darken or lighten and look wrong in one of the two modes).
func addPlasmaVariables(vars map[string]string) {
	bg, fg, accent := vars["background"], vars["foreground"], vars["accent"]
	if bg == "" || fg == "" {
		return
	}
	set := func(key, value string) {
		if vars[key] == "" {
			vars[key] = value
		}
	}
	set("kde_view", bg)
	set("kde_view_alt", mixHex(bg, fg, 0.04))
	set("kde_window", mixHex(bg, fg, 0.06))
	set("kde_window_alt", mixHex(bg, fg, 0.10))
	set("kde_button", mixHex(bg, fg, 0.12))
	set("kde_button_alt", mixHex(bg, fg, 0.16))
	set("kde_header", mixHex(bg, fg, 0.03))
	set("kde_header_inactive", bg)
	set("kde_tooltip", mixHex(bg, fg, 0.08))
	set("kde_fg_inactive", mixHex(fg, bg, 0.40))
	if accent != "" {
		set("kde_selection", accent)
		set("kde_selection_alt", mixHex(accent, bg, 0.30))
		set("kde_selection_fg", readableOn(accent))
		set("kde_selection_fg_inactive", mixHex(readableOn(accent), accent, 0.35))
		set("kde_hover", mixHex(accent, fg, 0.25))
	}
}
