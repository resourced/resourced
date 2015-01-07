package libstring

import (
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
