#!/bin/bash
# Switch every Plasma desktop to the Aether Plasma wallpaper plugin (adds
# "Open Aether Plasma" to the desktop right-click menu), keeping the current
# image. `activate.sh off` switches back to Plasma's Image wallpaper.
#
# Skipped quietly outside a running Plasma session.
set -euo pipefail

ID="io.github.davidboulay.aether-plasma-wallpaper"
MODE="${1:-on}"

qdbus=$(command -v qdbus6 || command -v qdbus || true)
if [ -z "$qdbus" ] || ! "$qdbus" org.kde.plasmashell /PlasmaShell >/dev/null 2>&1; then
    echo "Plasma is not running; the desktop menu item is added the next time Aether Plasma sets a wallpaper."
    exit 0
fi

if [ "$MODE" = off ]; then
    from="$ID" to="org.kde.image"
else
    from="org.kde.image" to="$ID"
fi

script=$(cat <<EOF
for (const d of desktops()) {
    if (d.wallpaperPlugin !== "$from") continue;
    d.currentConfigGroup = ["Wallpaper", "$from", "General"];
    const img = d.readConfig("Image");
    d.wallpaperPlugin = "$to";
    d.currentConfigGroup = ["Wallpaper", "$to", "General"];
    if (img) d.writeConfig("Image", img);
}
EOF
)
"$qdbus" org.kde.plasmashell /PlasmaShell org.kde.PlasmaShell.evaluateScript "$script" >/dev/null
echo "Desktop wallpaper plugin: $to"
