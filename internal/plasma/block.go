package plasma

import (
	"os"
	"path/filepath"
	"strings"

	"aether/internal/platform"
)

const (
	blockBegin = "# >>> aether-plasma (managed block, do not edit) >>>"
	blockEnd   = "# <<< aether-plasma <<<"
)

// stripManagedBlock removes an existing Aether Plasma block from text.
func stripManagedBlock(text string) string {
	start := strings.Index(text, blockBegin)
	if start < 0 {
		return text
	}
	end := strings.Index(text[start:], blockEnd)
	if end < 0 {
		return text
	}
	end += start + len(blockEnd)
	if end < len(text) && text[end] == '\n' {
		end++
	}
	before := strings.TrimRight(text[:start], "\n")
	after := text[end:]
	switch {
	case before == "":
		return strings.TrimLeft(after, "\n")
	case after == "":
		return before + "\n"
	default:
		return before + "\n" + after
	}
}

// ensureManagedBlock writes lines inside a marked block in path, at the top
// (prepend) or bottom. The file is created if missing and backed up once
// before Aether Plasma first changes it.
func ensureManagedBlock(path string, lines []string, prepend bool) error {
	block := blockBegin + "\n" + strings.Join(lines, "\n") + "\n" + blockEnd + "\n"

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	original := string(data)
	if strings.Contains(original, block) {
		return nil
	}
	rest := stripManagedBlock(original)

	var out string
	switch {
	case strings.TrimSpace(rest) == "":
		out = block
	case prepend:
		out = block + "\n" + rest
	default:
		out = strings.TrimRight(rest, "\n") + "\n\n" + block
	}

	if err := platform.EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	if original != "" {
		backupOnce(path)
	}
	return platform.WriteText(path, out)
}

// backupOnce keeps a copy of the user's file from before Aether Plasma first
// touched it, next to the file.
func backupOnce(path string) {
	backup := path + ".before-aether-plasma"
	if platform.FileExists(backup) || !platform.FileExists(path) {
		return
	}
	_ = platform.CopyFile(path, backup)
}
