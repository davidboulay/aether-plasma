package plasma

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/godbus/dbus/v5"

	"aether/internal/platform"
)

// WallpaperPluginID is the Aether Plasma wallpaper plugin
// (contrib/plasma-wallpaper). It shows the image Aether sets and adds
// "Open Aether Plasma" to the desktop right-click menu, which has no other
// extension point for third-party items.
const WallpaperPluginID = "io.github.davidboulay.aether-plasma-wallpaper"

// wallpaperPluginInstalled reports whether the plugin package is installed
// for this user or system-wide.
func wallpaperPluginInstalled() bool {
	dirs := []string{filepath.Join(dataHome(), "plasma", "wallpapers")}
	for _, d := range strings.Split(xdgDataDirs(), ":") {
		if d != "" {
			dirs = append(dirs, filepath.Join(d, "plasma", "wallpapers"))
		}
	}
	for _, d := range dirs {
		if platform.FileExists(filepath.Join(d, WallpaperPluginID, "metadata.json")) {
			return true
		}
	}
	return false
}

// wallpaperScript returns the Plasma desktop script that puts every desktop
// on the Aether wallpaper plugin with the given image URL.
func wallpaperScript(imageURL string) string {
	id, _ := json.Marshal(WallpaperPluginID)
	img, _ := json.Marshal(imageURL)
	return fmt.Sprintf(`for (const d of desktops()) {
    d.wallpaperPlugin = %s;
    d.currentConfigGroup = ["Wallpaper", %s, "General"];
    d.writeConfig("Image", %s);
}`, id, id, img)
}

// applyPluginWallpaper sets the image through the Aether wallpaper plugin.
func applyPluginWallpaper(abs string) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	obj := conn.Object("org.kde.plasmashell", "/PlasmaShell")
	return obj.Call("org.kde.PlasmaShell.evaluateScript", 0, wallpaperScript("file://"+abs)).Err
}
