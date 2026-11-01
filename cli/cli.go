package cli

import (
	"embed"
	"fmt"
	"os"
	"strings"
)

// Version is set at build time via ldflags.
var Version = "dev"

// Run dispatches CLI commands. Returns exit code.
func Run(args []string, templatesFS embed.FS) int {
	if len(args) == 0 {
		printUsage()
		return 1
	}

	cmd := args[0]
	switch cmd {
	// --- Existing commands ---
	case "--generate":
		return runGenerate(args[1:], templatesFS)
	case "--list-blueprints":
		return runListBlueprints(args[1:])
	case "--apply-blueprint":
		return runApplyBlueprint(args[1:], templatesFS)
	case "--import-blueprint":
		return runImportBlueprint(args[1:], templatesFS)
	case "--import-base16":
		return runImportBase16(args[1:], templatesFS)
	case "--import-colors-toml":
		return runImportColorsToml(args[1:], templatesFS)
	case "--handle-url":
		return runHandleURL(args[1:], templatesFS)

	// --- Color utilities ---
	case "--color-convert":
		return runColorConvert(args[1:])
	case "--contrast":
		return runContrast(args[1:])
	case "--adjust-color":
		return runAdjustColor(args[1:])
	case "--darken":
		return runDarken(args[1:])
	case "--lighten":
		return runLighten(args[1:])
	case "--color-info":
		return runColorInfo(args[1:])

	// --- Palette operations ---
	case "--extract-palette":
		return runExtractPalette(args[1:])
	case "--palette-from-color":
		return runPaletteFromColor(args[1:])
	case "--gradient":
		return runGradient(args[1:])
	case "--adjust-palette":
		return runAdjustPalette(args[1:])
	case "--palette-info":
		return runPaletteInfo(args[1:])

	// --- Blueprint extensions ---
	case "--show-blueprint":
		return runShowBlueprint(args[1:])
	case "--delete-blueprint":
		return runDeleteBlueprint(args[1:])
	case "--export-blueprint":
		return runExportBlueprint(args[1:])

	// --- Inspection ---
	case "--list-apps":
		return runListApps(args[1:], templatesFS)
	case "--list-modes":
		return runListModes(args[1:])
	case "--show-variables":
		return runShowVariables(args[1:], templatesFS)
	case "--current-theme":
		return runCurrentTheme(args[1:])
	case "--preview-template":
		return runPreviewTemplate(args[1:], templatesFS)

	// --- Wallpapers ---
	case "--search-wallhaven":
		return runSearchWallhaven(args[1:])
	case "--wallhaven-thumbs":
		return runWallhavenThumbs(args[1:])
	case "--wallhaven-download":
		return runDownloadWallhaven(args[1:])
	case "--wallhaven-info":
		return runWallhavenInfo(args[1:])
	case "--list-wallpapers":
		return runListWallpapers(args[1:])
	case "--random-wallpaper":
		return runRandomWallpaper(args[1:])

	// --- Favorites ---
	case "--list-favorites":
		return runListFavorites(args[1:])
	case "--toggle-favorite":
		return runToggleFavorite(args[1:])
	case "--is-favorite":
		return runIsFavorite(args[1:])

	// --- Meta ---
	case "upgrade":
		return runUpgrade(args[1:])
	case "--help", "-h":
		printUsage()
		return 0
	case "--version", "-v":
		fmt.Printf("aether-plasma %s\n", Version)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		return 1
	}
}

func printUsage() {
	fmt.Println(`Aether Plasma - Desktop Theme Generator for KDE Plasma

Usage:
  aether-plasma                             Launch GUI
  aether-plasma upgrade                     Upgrade to the latest release
  aether-plasma --help                      Show this help
  aether-plasma --version                   Show version

Theme generation:
  aether-plasma --generate <wallpaper> [options]   Extract colors and apply theme
    --light-mode                            Generate light variant
    --no-apply                              Render templates without activating
    --output <path>                         Output dir (defaults to ~/.config/aether-plasma/theme)
    --icon-theme automatic|<ID>              Use color-matched Yaru or an explicit installed theme ID
    --no-zed                                Skip Zed extension (default on)
    --no-vscode                             Skip VSCode integration (default on)
    --no-neovim                             Skip Neovim template (default on)
  aether-plasma --extract-palette <wallpaper>      Extract palette (read-only, no apply)

Blueprint management:
  aether-plasma --list-blueprints           List saved themes
  aether-plasma --show-blueprint <name>     Show blueprint details
  aether-plasma --apply-blueprint <name>    Apply a saved theme
  aether-plasma --delete-blueprint <name>   Delete a blueprint
  aether-plasma --export-blueprint <name>   Export blueprint to JSON file

Import commands:
  aether-plasma --import-blueprint <url|path>      Import blueprint from URL or file
  aether-plasma --import-base16 <url|file.yaml>    Import Base16 color scheme (URL or file)
    --wallpaper <url|path>                  Set wallpaper from URL or local path
    --light-mode                            Force light mode
  aether-plasma --import-colors-toml <url|file>    Import colors.toml color scheme (URL or file)
    --wallpaper <url|path>                  Set wallpaper from URL or local path
    --light-mode                            Force light mode
  aether-plasma --handle-url <aether://...> Handle a web link (confirms unless silent=true)
                                              aether://apply?external_theme=URL
                                              aether://apply?colors=URL
                                              aether://apply?wallpaper=URL
                                              aether://apply?mode=light|dark
                                              aether://apply?as_omarchy_theme=NAME

Color utilities:
  aether-plasma --color-info <hex>          Show all color representations
  aether-plasma --color-convert <hex> [--to fmt]   Convert color (rgb|hsl|oklab|oklch)
  aether-plasma --contrast <color1> <color2>       WCAG contrast ratio + grading
  aether-plasma --adjust-color <hex> [adjustments] Apply adjustment pipeline
  aether-plasma --darken <hex> <percent>    Darken a color
  aether-plasma --lighten <hex> <percent>   Lighten a color

Palette operations:
  aether-plasma --palette-from-color <hex>  Generate 16-color palette from one color
  aether-plasma --gradient <start> <end>    Generate color gradient
  aether-plasma --adjust-palette <colors> [adj]    Adjust entire palette
  aether-plasma --palette-info <colors|blueprint>  Analyze palette

Inspection:
  aether-plasma --list-apps                 List supported applications
  aether-plasma --list-modes                List extraction modes
  aether-plasma --show-variables <source>   Show template variables
  aether-plasma --current-theme             Show currently applied theme
  aether-plasma --preview-template <app> <source>  Preview rendered template

Wallpapers:
  aether-plasma --search-wallhaven <query>  Search wallhaven.cc (metadata only)
  aether-plasma --wallhaven-thumbs <query>  Search + download small thumbnails to cache
  aether-plasma --wallhaven-download <id|url>      Download full wallpaper by id or page URL
  aether-plasma --wallhaven-info <id|url>   Show metadata for a single wallpaper
  aether-plasma --list-wallpapers [--with-previews] List local wallpapers (cached 800px previews when flag set)
  aether-plasma --random-wallpaper          Pick random local wallpaper

  Wallhaven search flags (apply to --search-wallhaven and --wallhaven-thumbs):
    --page N                                Starting page number (default 1)
    --pages N                               (--wallhaven-thumbs only) Fetch N consecutive pages concurrently
    --categories 111                        Bitmask: General/Anime/People
    --purity 100                            Bitmask: SFW/Sketchy/NSFW
    --sorting date_added|relevance|random|views|favorites|toplist
    --order desc|asc
    --at-least 1920x1080                    Minimum resolution
    --colors HEX                            Filter by wallhaven palette color

Favorites:
  aether-plasma --list-favorites            List all favorites
  aether-plasma --toggle-favorite <path>    Add/remove favorite
  aether-plasma --is-favorite <path>        Check if favorited

Remote control (requires running Aether GUI):
  aether-plasma status                      Show editor state
  aether-plasma extract <wallpaper>         Load wallpaper into editor
  aether-plasma set-color <index> <hex>     Change palette color
  aether-plasma adjust [--vibrance X] ...   Adjust palette sliders
  aether-plasma set-mode <mode>             Change extraction mode
  aether-plasma apply                       Apply current theme
  aether-plasma toggle-light-mode           Toggle light/dark mode
  aether-plasma load-blueprint <name>       Load blueprint into editor
  aether-plasma apply-blueprint <name>      Apply blueprint directly
  aether-plasma set-wallpaper <path>        Set wallpaper path

GUI options:
  aether-plasma --tab <name>                Open GUI with specific tab

Global options:
  --json                                    Machine-readable JSON output
  --light-mode                              Use light mode
  --extract-mode <mode>                     Extraction mode (see --list-modes)`)
}

func parseFlag(args []string, flag string) (string, []string) {
	for i, arg := range args {
		if arg == flag && i+1 < len(args) {
			remaining := make([]string, 0, len(args)-2)
			remaining = append(remaining, args[:i]...)
			remaining = append(remaining, args[i+2:]...)
			return args[i+1], remaining
		}
	}
	return "", args
}

func hasFlag(args []string, flag string) (bool, []string) {
	for i, arg := range args {
		if arg == flag {
			remaining := make([]string, 0, len(args)-1)
			remaining = append(remaining, args[:i]...)
			remaining = append(remaining, args[i+1:]...)
			return true, remaining
		}
	}
	return false, args
}

// validModes are the allowed extraction modes.
var validModes = map[string]bool{
	"normal": true, "monochromatic": true, "analogous": true,
	"pastel": true, "material": true, "colorful": true,
	"muted": true, "bright": true,
	"complementary": true, "triadic": true,
	"split-complementary": true, "tetradic": true,
	"fire": true, "ocean": true, "forest": true,
	"earthtone": true, "neon": true, "sunset": true, "vaporwave": true,
	"midnight": true, "aurora": true,
	"high-contrast": true, "duotone": true,
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return home + path[1:]
		}
	}
	return path
}
