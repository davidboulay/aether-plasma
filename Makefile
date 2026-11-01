.PHONY: build dev test clean install uninstall update install-omarchy-plugins

WEBKIT_TAGS := $(shell pkg-config --exists webkit2gtk-4.0 2>/dev/null || echo "-tags webkit2_41")
PREFIX ?= $(HOME)/.local
BINDIR := $(PREFIX)/bin
APPDIR := $(PREFIX)/share/applications
ICONDIR := $(PREFIX)/share/icons/hicolor/512x512/apps
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X 'aether/internal/update.packageUpdateCommand=make -C $(CURDIR) update' -X 'aether/cli.Version=$(VERSION)'

build:
	wails build $(WEBKIT_TAGS) -ldflags "$(LDFLAGS)"

# User-local install (no sudo): ~/.local/bin/aether-plasma + desktop entry.
install: build
	install -Dm755 build/bin/aether-plasma $(BINDIR)/aether-plasma
	install -Dm644 dev.davidboulay.aether-plasma.desktop $(APPDIR)/dev.davidboulay.aether-plasma.desktop
	install -Dm644 assets/aether-icon-512.png $(ICONDIR)/aether-plasma.png
	-update-desktop-database $(APPDIR) 2>/dev/null
	-gtk-update-icon-cache -f $(PREFIX)/share/icons/hicolor 2>/dev/null
	-kbuildsycoca6 2>/dev/null

uninstall:
	rm -f $(BINDIR)/aether-plasma $(APPDIR)/dev.davidboulay.aether-plasma.desktop $(ICONDIR)/aether-plasma.png
	-update-desktop-database $(APPDIR) 2>/dev/null

update:
	git -C $(CURDIR) pull --ff-only
	$(MAKE) -C $(CURDIR) install

install-omarchy-plugins:
	AETHER_PLUGIN_SOURCE_DIR="$(CURDIR)/contrib/quickshell" contrib/quickshell/install.sh

dev:
	wails dev $(WEBKIT_TAGS)

test:
	go test ./internal/... ./cli/...
	go test -ldflags "-X aether/internal/update.packageUpdateCommand=omarchy-update" ./internal/update ./cli

clean:
	rm -rf build/bin
