/*
 * Aether Plasma wallpaper.
 *
 * Plasma's desktop right-click menu has no extension point for third-party
 * items; the only part a user-installed package can add to is the wallpaper
 * plugin's own actions. This plugin shows the image Aether Plasma sets and
 * adds "Open Aether Plasma" there.
 *
 * SPDX-License-Identifier: MIT
 */

import QtQuick
import org.kde.plasma.core as PlasmaCore
import org.kde.plasma.plasmoid
import org.kde.plasma.plasma5support as P5Support

WallpaperItem {
    id: root

    readonly property string imageUrl: {
        const value = String(root.configuration.Image || "");
        return value === "" || value.indexOf("://") !== -1 ? value : "file://" + value;
    }

    contextualActions: [
        PlasmaCore.Action {
            text: "Open Aether Plasma"
            icon.name: "aether-plasma"
            onTriggered: launcher.run("kstart --application aether-plasma")
        },
        PlasmaCore.Action {
            text: "Open Wallpaper Image"
            icon.name: "document-open"
            enabled: root.imageUrl !== ""
            onTriggered: Qt.openUrlExternally(root.imageUrl)
        }
    ]

    P5Support.DataSource {
        id: launcher
        engine: "executable"
        connectedSources: []
        onNewData: (sourceName, data) => disconnectSource(sourceName)
        function run(command) {
            connectSource(command);
        }
    }

    Rectangle {
        anchors.fill: parent
        color: root.configuration.Color
    }

    Image {
        anchors.fill: parent
        source: root.imageUrl
        fillMode: root.configuration.FillMode
        sourceSize.width: root.width * Screen.devicePixelRatio
        sourceSize.height: root.height * Screen.devicePixelRatio
        asynchronous: true
        cache: false
        smooth: true
        // Hold the old frame until the new one is decoded, so a theme switch
        // never flashes the background color.
        retainWhileLoading: true
    }
}
