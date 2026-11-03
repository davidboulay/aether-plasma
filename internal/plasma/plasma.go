// Package plasma applies a generated Aether theme to a KDE Plasma session:
// color scheme, wallpaper (desktop + lock screen), icon theme, and the
// terminal/TUI apps that can include a generated file (Konsole, Ghostty,
// Kitty, Alacritty, Foot, btop).
//
// Every integration is a layer on top of the user's own config: Aether Plasma
// only writes its own files (color scheme, Konsole scheme, btop theme) and a
// marked include block in app configs, so a Plasma or app update never
// overwrites it and removing the block restores the original behavior.
package plasma

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"aether/internal/platform"
)

const (
	// SchemeName is the display name of the generated KDE color scheme.
	SchemeName = "Aether Plasma"
	schemeA    = "AetherPlasmaA"
	schemeB    = "AetherPlasmaB"

	konsoleScheme  = "AetherPlasma"
	konsoleProfile = "Aether Plasma.profile"
	btopThemeName  = "aether-plasma"
)

// errNotInstalled marks an app target that is skipped because the app is
// not installed; it is not reported as a warning.
var errNotInstalled = errors.New("not installed")

// Options describes one theme application.
type Options struct {
	ThemeDir  string            // directory with the rendered templates
	Wallpaper string            // absolute wallpaper path, may be empty
	LightMode bool              // light variant
	Include   func(string) bool // whether an optional app target is enabled
}

// Result reports what was applied and what was skipped.
type Result struct {
	Applied  []string `json:"applied"`
	Warnings []string `json:"warnings"`
}

func (r *Result) ok(what string) { r.Applied = append(r.Applied, what) }
func (r *Result) warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("aether-plasma: %s", msg)
	r.Warnings = append(r.Warnings, msg)
}

// Detected reports whether the current session is KDE Plasma with the
// plasma-apply-* tools available. AETHER_PLASMA=0 disables it, =1 forces it.
func Detected() bool {
	switch os.Getenv("AETHER_PLASMA") {
	case "0":
		return false
	case "1":
		return true
	}
	desktop := strings.ToUpper(os.Getenv("XDG_CURRENT_DESKTOP") + ":" + os.Getenv("XDG_SESSION_DESKTOP"))
	if !strings.Contains(desktop, "KDE") && !strings.Contains(desktop, "PLASMA") {
		return false
	}
	return platform.CommandExists("plasma-apply-colorscheme")
}

// Apply activates the rendered theme in the running Plasma session. Failures
// of individual integrations are collected as warnings; only a failure of the
// color scheme itself is returned as an error.
func Apply(opts Options) (*Result, error) {
	res := &Result{}
	include := opts.Include
	if include == nil {
		include = func(string) bool { return true }
	}

	if err := applyColorScheme(opts.ThemeDir); err != nil {
		return res, fmt.Errorf("apply Plasma color scheme: %w", err)
	}
	res.ok("color scheme")

	switch changed, err := ensureSchemeFollowingStyle(); {
	case err != nil:
		res.warn("panel and menus: %v", err)
	case changed:
		res.ok("Plasma style (Breeze, follows the scheme)")
	}

	if opts.Wallpaper != "" {
		if err := applyWallpaper(opts.Wallpaper); err != nil {
			res.warn("wallpaper: %v", err)
		} else {
			res.ok("wallpaper")
		}
	}

	if include("icons") {
		if name, err := applyIconTheme(opts.ThemeDir); err != nil {
			res.warn("icons: %v", err)
		} else if name != "" {
			res.ok("icons (" + name + ")")
		}
	}

	type integration struct {
		app string
		fn  func(string) error
	}
	for _, it := range []integration{
		{"konsole", applyKonsole},
		{"ghostty", applyGhostty},
		{"kitty", applyKitty},
		{"alacritty", applyAlacritty},
		{"foot", applyFoot},
		{"btop", applyBtop},
	} {
		if !include(it.app) {
			continue
		}
		if err := it.fn(opts.ThemeDir); err != nil {
			if !errors.Is(err, errNotInstalled) {
				res.warn("%s: %v", it.app, err)
			}
			continue
		}
		res.ok(it.app)
	}
	return res, nil
}

// Restore switches back to the color scheme that was active before Aether
// Plasma first applied one.
func Restore() error {
	state := loadState()
	if state.PreviousScheme == "" {
		return nil
	}
	if _, err := platform.RunSync("plasma-apply-colorscheme", state.PreviousScheme); err != nil {
		return fmt.Errorf("restore color scheme %q: %w", state.PreviousScheme, err)
	}
	for _, name := range []string{schemeA, schemeB} {
		_ = os.Remove(filepath.Join(colorSchemesDir(), name+".colors"))
	}
	if state.PreviousStyle != "" {
		if _, err := platform.RunSync("plasma-apply-desktoptheme", state.PreviousStyle); err != nil {
			log.Printf("aether-plasma: restore Plasma style %q: %v", state.PreviousStyle, err)
		}
	}
	state.PreviousScheme = ""
	state.PreviousStyle = ""
	return saveState(state)
}

// ---------------------------------------------------------------------------
// paths and state

func homeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func dataHome() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return xdg
	}
	return filepath.Join(homeDir(), ".local", "share")
}

func configHome() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	return filepath.Join(homeDir(), ".config")
}

func colorSchemesDir() string { return filepath.Join(dataHome(), "color-schemes") }
func konsoleDir() string      { return filepath.Join(dataHome(), "konsole") }

type state struct {
	PreviousScheme string `json:"previousScheme"`
	PreviousStyle  string `json:"previousStyle,omitempty"`
}

func statePath() string { return filepath.Join(platform.ConfigDir(), "plasma-state.json") }

func loadState() state {
	s, err := platform.ReadJSON[state](statePath())
	if err != nil {
		return state{}
	}
	return s
}

func saveState(s state) error { return platform.WriteJSON(statePath(), s) }

func readConfig(file string, groups []string, key string) string {
	args := []string{"--file", file}
	for _, g := range groups {
		args = append(args, "--group", g)
	}
	args = append(args, "--key", key)
	out, err := platform.RunSync("kreadconfig6", args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

func writeConfig(file string, groups []string, key, value string) error {
	args := []string{"--file", file}
	for _, g := range groups {
		args = append(args, "--group", g)
	}
	args = append(args, "--key", key, value)
	_, err := platform.RunSync("kwriteconfig6", args...)
	return err
}

// ---------------------------------------------------------------------------
// color scheme

// applyColorScheme installs kde.colors under an alternating file name.
// plasma-apply-colorscheme ignores a request for the scheme that is already
// active, so re-applying a changed palette needs a different scheme name.
func applyColorScheme(themeDir string) error {
	content, err := os.ReadFile(filepath.Join(themeDir, "kde.colors"))
	if err != nil {
		return fmt.Errorf("read generated scheme: %w", err)
	}

	current := readConfig("kdeglobals", []string{"General"}, "ColorScheme")
	next, stale := schemeA, schemeB
	if current == schemeA {
		next, stale = schemeB, schemeA
	}
	if current != schemeA && current != schemeB && current != "" {
		s := loadState()
		s.PreviousScheme = current
		if err := saveState(s); err != nil {
			log.Printf("aether-plasma: could not record previous scheme: %v", err)
		}
	}

	body := strings.Replace(string(content), "ColorScheme=AetherPlasma\n", "ColorScheme="+next+"\n", 1)
	if err := platform.EnsureDir(colorSchemesDir()); err != nil {
		return err
	}
	target := filepath.Join(colorSchemesDir(), next+".colors")
	if err := platform.WriteText(target, body); err != nil {
		return err
	}
	if _, err := platform.RunSync("plasma-apply-colorscheme", next); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(colorSchemesDir(), stale+".colors")); err != nil && !os.IsNotExist(err) {
		log.Printf("aether-plasma: could not remove stale scheme: %v", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Plasma style (panel, menus, widgets)

// fixedBreezeStyles have hard-coded colors and ignore the color scheme;
// "default" (Breeze) draws the panel and menus from the active scheme.
var fixedBreezeStyles = map[string]bool{"breeze-dark": true, "breeze-light": true}

// ensureSchemeFollowingStyle switches Breeze Dark/Light to Breeze so the
// panel and menus follow the generated scheme. It remembers the previous
// style for Restore. Third-party styles are left alone with a warning.
func ensureSchemeFollowingStyle() (bool, error) {
	current := readConfig("plasmarc", []string{"Theme"}, "name")
	if current == "" || current == "default" {
		return false, nil
	}
	if !fixedBreezeStyles[current] {
		return false, fmt.Errorf("Plasma style %q may use its own colors; choose Breeze in System Settings → Plasma Style for the panel to follow the theme", current)
	}
	s := loadState()
	if s.PreviousStyle == "" {
		s.PreviousStyle = current
		if err := saveState(s); err != nil {
			log.Printf("aether-plasma: could not record previous style: %v", err)
		}
	}
	if _, err := platform.RunSync("plasma-apply-desktoptheme", "default"); err != nil {
		return false, err
	}
	return true, nil
}

// ---------------------------------------------------------------------------
// wallpaper

// SetWallpaper changes the desktop and lock screen wallpaper without touching
// the color scheme.
func SetWallpaper(path string) error { return applyWallpaper(path) }

func applyWallpaper(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); err != nil {
		return err
	}
	// With the Aether wallpaper plugin installed, set the image through it so
	// the desktop keeps its "Open Aether Plasma" menu item;
	// plasma-apply-wallpaperimage would switch it back to the stock plugin.
	if !wallpaperPluginInstalled() || applyPluginWallpaper(abs) != nil {
		if _, err := platform.RunSync("plasma-apply-wallpaperimage", abs); err != nil {
			return err
		}
	}
	// Lock screen uses the same image.
	groups := []string{"Greeter", "Wallpaper", "org.kde.image", "General"}
	if err := writeConfig("kscreenlockerrc", groups, "Image", "file://"+abs); err != nil {
		return fmt.Errorf("lock screen: %w", err)
	}
	_ = writeConfig("kscreenlockerrc", groups, "PreviewImage", "file://"+abs)
	return nil
}

// ---------------------------------------------------------------------------
// icons

func xdgDataDirs() string {
	if d := os.Getenv("XDG_DATA_DIRS"); d != "" {
		return d
	}
	return "/usr/local/share:/usr/share"
}

func iconThemeInstalled(name string) bool {
	dirs := []string{
		filepath.Join(dataHome(), "icons"),
		filepath.Join(homeDir(), ".icons"),
	}
	for _, d := range strings.Split(xdgDataDirs(), ":") {
		if d != "" {
			dirs = append(dirs, filepath.Join(d, "icons"))
		}
	}
	for _, d := range dirs {
		if platform.FileExists(filepath.Join(d, name, "index.theme")) {
			return true
		}
	}
	return false
}

// applyIconTheme switches to the icon theme named in icons.theme when it is
// installed. A missing theme (e.g. the automatic Yaru variant) is skipped
// silently so the user's current Plasma icons stay.
func applyIconTheme(themeDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(themeDir, "icons.theme"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	name := strings.TrimSpace(strings.SplitN(string(data), "\n", 2)[0])
	if name == "" || !iconThemeInstalled(name) {
		return "", nil
	}
	if readConfig("kdeglobals", []string{"Icons"}, "Theme") == name {
		return "", nil
	}
	tool := "/usr/lib/plasma-changeicons"
	if !platform.FileExists(tool) {
		return "", fmt.Errorf("plasma-changeicons not found")
	}
	if _, err := platform.RunSync(tool, name); err != nil {
		return "", err
	}
	return name, nil
}

// ---------------------------------------------------------------------------
// Konsole

func applyKonsole(themeDir string) error {
	if !platform.CommandExists("konsole") {
		return errNotInstalled
	}
	content, err := os.ReadFile(filepath.Join(themeDir, "konsole.colorscheme"))
	if err != nil {
		return err
	}
	if err := platform.EnsureDir(konsoleDir()); err != nil {
		return err
	}
	if err := platform.WriteText(filepath.Join(konsoleDir(), konsoleScheme+".colorscheme"), string(content)); err != nil {
		return err
	}

	profile := readConfig("konsolerc", []string{"Desktop Entry"}, "DefaultProfile")
	if profile == "" {
		profile = konsoleProfile
		path := filepath.Join(konsoleDir(), profile)
		if !platform.FileExists(path) {
			body := "[General]\nName=Aether Plasma\nParent=FALLBACK/\n\n[Appearance]\nColorScheme=" + konsoleScheme + "\n"
			if err := platform.WriteText(path, body); err != nil {
				return err
			}
		}
		if err := writeConfig("konsolerc", []string{"Desktop Entry"}, "DefaultProfile", profile); err != nil {
			return err
		}
	}
	profilePath := filepath.Join(konsoleDir(), profile)
	if !platform.FileExists(profilePath) {
		return fmt.Errorf("default profile %q not found in %s", profile, konsoleDir())
	}
	return writeConfig(profilePath, []string{"Appearance"}, "ColorScheme", konsoleScheme)
}

// ---------------------------------------------------------------------------
// include-based terminals

func firstExisting(paths ...string) string {
	for _, p := range paths {
		if platform.FileExists(p) {
			return p
		}
	}
	return ""
}

func applyGhostty(themeDir string) error {
	if !platform.CommandExists("ghostty") {
		return errNotInstalled
	}
	dir := filepath.Join(configHome(), "ghostty")
	cfg := firstExisting(filepath.Join(dir, "config.ghostty"), filepath.Join(dir, "config"))
	if cfg == "" {
		cfg = filepath.Join(dir, "config.ghostty")
	}
	line := fmt.Sprintf("config-file = ?%q", filepath.Join(themeDir, "ghostty.conf"))
	if err := ensureManagedBlock(cfg, []string{line}, false); err != nil {
		return err
	}
	signalProcess("ghostty", "USR2") // Ghostty ≥ 1.2 reloads config on SIGUSR2
	return nil
}

func applyKitty(themeDir string) error {
	if !platform.CommandExists("kitty") {
		return errNotInstalled
	}
	cfg := filepath.Join(configHome(), "kitty", "kitty.conf")
	line := "include " + filepath.Join(themeDir, "kitty.conf")
	if err := ensureManagedBlock(cfg, []string{line}, false); err != nil {
		return err
	}
	signalProcess("kitty", "USR1") // kitty reloads config on SIGUSR1
	return nil
}

func applyFoot(themeDir string) error {
	if !platform.CommandExists("foot") {
		return errNotInstalled
	}
	dir := filepath.Join(configHome(), "foot")
	// foot only accepts include= in the main section, so prepend.
	line := "include=" + filepath.Join(themeDir, "foot.ini")
	return ensureManagedBlock(filepath.Join(dir, "foot.ini"), []string{line}, true)
}

func applyAlacritty(themeDir string) error {
	if !platform.CommandExists("alacritty") {
		return errNotInstalled
	}
	cfg := filepath.Join(configHome(), "alacritty", "alacritty.toml")
	themeFile := filepath.Join(themeDir, "alacritty.toml")
	if data, err := os.ReadFile(cfg); err == nil {
		text := string(stripManagedBlock(string(data)))
		if strings.Contains(text, themeFile) {
			return nil // user imports it already
		}
		if strings.Contains(text, "[general]") {
			return fmt.Errorf("%s has a [general] table; add %q to general.import yourself", cfg, themeFile)
		}
	}
	// TOML dotted keys must precede tables, so prepend. Alacritty reloads
	// its config live when the imported file changes.
	line := fmt.Sprintf("general.import = [%q]", themeFile)
	return ensureManagedBlock(cfg, []string{line}, true)
}

// ---------------------------------------------------------------------------
// btop

func applyBtop(themeDir string) error {
	if !platform.CommandExists("btop") {
		return errNotInstalled
	}
	content, err := os.ReadFile(filepath.Join(themeDir, "btop.theme"))
	if err != nil {
		return err
	}
	dir := filepath.Join(configHome(), "btop")
	if err := platform.EnsureDir(filepath.Join(dir, "themes")); err != nil {
		return err
	}
	if err := platform.WriteText(filepath.Join(dir, "themes", btopThemeName+".theme"), string(content)); err != nil {
		return err
	}

	cfg := filepath.Join(dir, "btop.conf")
	want := fmt.Sprintf("color_theme = %q", btopThemeName)
	data, err := os.ReadFile(cfg)
	switch {
	case os.IsNotExist(err):
		if err := platform.WriteText(cfg, want+"\n"); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		lines := strings.Split(string(data), "\n")
		found := false
		for i, l := range lines {
			if strings.HasPrefix(strings.TrimSpace(l), "color_theme") {
				if strings.TrimSpace(l) != want {
					backupOnce(cfg)
				}
				lines[i] = want
				found = true
			}
		}
		if !found {
			backupOnce(cfg)
			lines = append([]string{want}, lines...)
		}
		if err := platform.WriteText(cfg, strings.Join(lines, "\n")); err != nil {
			return err
		}
	}
	signalProcess("btop", "USR2") // btop reloads config on SIGUSR2
	return nil
}

// ---------------------------------------------------------------------------
// helpers

func signalProcess(name, sig string) {
	if !platform.CommandExists("pkill") {
		return
	}
	// pkill exits 1 when nothing matched; that is fine.
	_, _ = platform.RunSync("pkill", "-"+sig, "-x", name)
}
