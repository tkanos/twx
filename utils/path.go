package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ExpandTilde converts ~/path/to/file to /home/user/path/to/file.
func ExpandTilde(dir string) string {
	if strings.HasPrefix(dir, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, dir[2:])
		}
	}

	return dir
}
