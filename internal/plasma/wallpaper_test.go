package plasma

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWallpaperScriptQuotesValues(t *testing.T) {
	s := wallpaperScript(`file:///home/me/it's "here".jpg`)
	for _, want := range []string{
		`d.wallpaperPlugin = "` + WallpaperPluginID + `";`,
		`["Wallpaper", "` + WallpaperPluginID + `", "General"]`,
		`d.writeConfig("Image", "file:///home/me/it's \"here\".jpg");`,
	} {
		if !strings.Contains(s, want) {
			t.Errorf("script missing %q:\n%s", want, s)
		}
	}
}

func TestWallpaperPluginInstalled(t *testing.T) {
	data := t.TempDir()
	t.Setenv("XDG_DATA_HOME", data)
	t.Setenv("XDG_DATA_DIRS", t.TempDir())
	if wallpaperPluginInstalled() {
		t.Fatal("reported installed with no package")
	}
	dir := filepath.Join(data, "plasma", "wallpapers", WallpaperPluginID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "metadata.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !wallpaperPluginInstalled() {
		t.Fatal("not detected in XDG_DATA_HOME")
	}
}
