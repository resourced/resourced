// Package libstring provides string related library functions.
package libstring

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"os/user"
	"strings"
)

// ExpandTilde is a convenience function that expands ~ to full path.
func ExpandTilde(path string) string {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	if path[:2] == "~/" {
		path = strings.Replace(path, "~", homeDir, 1)
	}
	return path
}

// ExpandTilde is a convenience function that expands both ~ and $ENV.
func ExpandTildeAndEnv(path string) string {
	path = ExpandTilde(path)
	return os.ExpandEnv(path)
}

// GeneratePassword returns password.
// size determines length of initial seed bytes.
func GeneratePassword(size int) (string, error) {
	// Force minimum size to 32
	if size < 32 {
		size = 32
	}

	rb := make([]byte, size)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}
