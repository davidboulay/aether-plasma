package template

import (
	"os"
	"regexp"
	"testing"
)

func TestPlasmaVariablesDerivedInBothModes(t *testing.T) {
	for _, light := range []bool{false, true} {
		roles := ColorRoles{Background: "#1e1e2e", Foreground: "#cdd6f4", Blue: "#89b4fa"}
		if light {
			roles.Background, roles.Foreground = "#eff1f5", "#4c4f69"
		}
		vars := BuildVariables(roles, light, nil)
		for _, key := range plasmaVariableKeys {
			if !regexp.MustCompile(`^#[0-9a-f]{6}$`).MatchString(vars[key]) {
				t.Errorf("light=%v %s = %q", light, key, vars[key])
			}
		}
		if vars["kde_view"] != roles.Background {
			t.Errorf("kde_view = %s, want background", vars["kde_view"])
		}
	}
}

func TestPlasmaVariablesFollowOverrides(t *testing.T) {
	roles := ColorRoles{Background: "#000000", Foreground: "#ffffff", Blue: "#0000ff"}
	vars := BuildVariables(roles, false, nil)
	vars["background"] = "#202020"
	RecomputeDerived(vars, map[string]string{"background": "#202020", "kde_tooltip": "#123456"})
	vars["kde_tooltip"] = "#123456"
	RecomputeDerived(vars, map[string]string{"background": "#202020", "kde_tooltip": "#123456"})
	if vars["kde_view"] != "#202020" {
		t.Fatalf("kde_view did not follow background override: %s", vars["kde_view"])
	}
	if vars["kde_tooltip"] != "#123456" {
		t.Fatalf("explicit kde_tooltip lost: %s", vars["kde_tooltip"])
	}
}

func TestKDETemplatesRenderCompletely(t *testing.T) {
	roles := ColorRoles{
		Background: "#1e1e2e", Foreground: "#cdd6f4", Black: "#45475a", Red: "#f38ba8",
		Green: "#a6e3a1", Yellow: "#f9e2af", Blue: "#89b4fa", Magenta: "#f5c2e7",
		Cyan: "#94e2d5", White: "#bac2de", BrightBlack: "#585b70", BrightRed: "#f38ba8",
		BrightGreen: "#a6e3a1", BrightYellow: "#f9e2af", BrightBlue: "#89b4fa",
		BrightMagenta: "#f5c2e7", BrightCyan: "#94e2d5", BrightWhite: "#a6adc8",
	}
	vars := BuildVariables(roles, false, nil)
	for _, name := range []string{"kde.colors", "konsole.colorscheme"} {
		raw, err := os.ReadFile("../../templates/" + name)
		if err != nil {
			t.Fatal(err)
		}
		out := ProcessTemplate(string(raw), vars)
		if m := regexp.MustCompile(`\{[a-z_]+(\.[a-z]+)?\}`).FindString(out); m != "" {
			t.Errorf("%s: unrendered variable %s", name, m)
		}
		if !regexp.MustCompile(`BackgroundNormal=30,30,46|Color=30,30,46`).MatchString(out) {
			t.Errorf("%s: background not rendered as r,g,b", name)
		}
	}
}
