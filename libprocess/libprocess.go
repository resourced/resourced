package libprocess

import (
	"os"
	"os/exec"
	"strings"
)

func NewCmd(command string) *exec.Cmd {
	wd, _ := os.Getwd()

	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := exec.Command(head, parts...)
	cmd.Dir = wd
	cmd.Env = os.Environ()

	return cmd
}
