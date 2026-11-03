/*
 * SPDX-License-Identifier: MIT
 */

import QtQuick
import QtQuick.Controls
import QtQuick.Dialogs
import QtQuick.Layouts
import org.kde.kirigami as Kirigami

Kirigami.FormLayout {
    id: root
    twinFormLayouts: parentLayout

    property string cfg_Image
    property int cfg_FillMode
    property alias cfg_Color: colorButton.color

    Label {
        Layout.maximumWidth: Kirigami.Units.gridUnit * 24
        wrapMode: Text.Wrap
        text: "Aether Plasma sets this image when it applies a theme. Right-click the desktop to open Aether Plasma."
    }

    Button {
        Kirigami.FormData.label: "Image:"
        Accessible.name: "Select an image"
        ToolTip.text: root.cfg_Image
        ToolTip.visible: hovered && root.cfg_Image !== ""
        icon.source: root.cfg_Image
        icon.width: Kirigami.Units.gridUnit * 8
        icon.height: Kirigami.Units.gridUnit * 5
        text: root.cfg_Image === "" ? "Select Image…" : ""
        onClicked: fileDialog.open()
    }

    ComboBox {
        Kirigami.FormData.label: "Positioning:"
        textRole: "text"
        valueRole: "value"
        model: [
            { text: "Scaled and Cropped", value: Image.PreserveAspectCrop },
            { text: "Scaled, Keep Proportions", value: Image.PreserveAspectFit },
            { text: "Scaled", value: Image.Stretch },
            { text: "Centered", value: Image.Pad }
        ]
        Component.onCompleted: currentIndex = Math.max(0, indexOfValue(root.cfg_FillMode))
        onActivated: root.cfg_FillMode = currentValue
    }

    Button {
        id: colorButton
        property color color: "black"
        Kirigami.FormData.label: "Background:"
        Accessible.name: "Background color"
        contentItem: Rectangle {
            implicitWidth: Kirigami.Units.gridUnit * 3
            implicitHeight: Kirigami.Units.gridUnit
            color: colorButton.color
            border.color: Kirigami.Theme.textColor
            border.width: 1
        }
        onClicked: colorDialog.open()
    }

    FileDialog {
        id: fileDialog
        fileMode: FileDialog.OpenFile
        options: FileDialog.ReadOnly
        nameFilters: ["Images (*.png *.jpg *.jpeg *.webp *.avif *.jxl *.bmp)"]
        title: "Open Image"
        onAccepted: root.cfg_Image = selectedFile
    }

    ColorDialog {
        id: colorDialog
        selectedColor: colorButton.color
        onAccepted: colorButton.color = selectedColor
    }
}
