package plasma

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureManagedBlockAppendsOnceAndBacksUp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	if err := os.WriteFile(path, []byte("font-size = 12\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := ensureManagedBlock(path, []string{"config-file = ?\"/x/ghostty.conf\""}, false); err != nil {
			t.Fatal(err)
		}
	}
	got, _ := os.ReadFile(path)
	text := string(got)
	if strings.Count(text, blockBegin) != 1 {
		t.Fatalf("block written %d times:\n%s", strings.Count(text, blockBegin), text)
	}
	if !strings.HasPrefix(text, "font-size = 12\n") {
		t.Fatalf("user content not kept first:\n%s", text)
	}
	backup, err := os.ReadFile(path + ".before-aether-plasma")
	if err != nil || string(backup) != "font-size = 12\n" {
		t.Fatalf("backup = %q, %v", backup, err)
	}
}

func TestEnsureManagedBlockReplacesAndPrepends(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alacritty.toml")
	if err := os.WriteFile(path, []byte("[font]\nsize = 11\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ensureManagedBlock(path, []string{"old"}, true); err != nil {
		t.Fatal(err)
	}
	if err := ensureManagedBlock(path, []string{"new"}, true); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	text := string(got)
	if strings.Contains(text, "old") || !strings.HasPrefix(text, blockBegin+"\nnew\n") {
		t.Fatalf("unexpected content:\n%s", text)
	}
	if !strings.HasSuffix(text, "[font]\nsize = 11\n") {
		t.Fatalf("user content lost:\n%s", text)
	}
}

func TestEnsureManagedBlockCreatesMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kitty", "kitty.conf")
	if err := ensureManagedBlock(path, []string{"include /x"}, false); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != blockBegin+"\ninclude /x\n"+blockEnd+"\n" {
		t.Fatalf("got %q", got)
	}
	if _, err := os.Stat(path + ".before-aether-plasma"); !os.IsNotExist(err) {
		t.Fatalf("backup created for a new file")
	}
}

func TestStripManagedBlock(t *testing.T) {
	in := "a\n\n" + blockBegin + "\nx\n" + blockEnd + "\nb\n"
	if got := stripManagedBlock(in); got != "a\nb\n" {
		t.Fatalf("got %q", got)
	}
}

func TestDetectedHonorsOverride(t *testing.T) {
	t.Setenv("AETHER_PLASMA", "0")
	t.Setenv("XDG_CURRENT_DESKTOP", "KDE")
	if Detected() {
		t.Fatal("AETHER_PLASMA=0 must disable detection")
	}
	t.Setenv("AETHER_PLASMA", "")
	t.Setenv("XDG_CURRENT_DESKTOP", "Hyprland")
	t.Setenv("XDG_SESSION_DESKTOP", "")
	if Detected() {
		t.Fatal("Hyprland session detected as Plasma")
	}
}
