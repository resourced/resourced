// +build darwin
package readers

import (
	"encoding/json"
	"os/exec"
	"strings"
)

func init() {
	Register("Uname", NewUname)
}

func NewUname() IReader {
	u := &Uname{}
	u.Data = make(map[string]interface{})
	return u
}

// Uname is a reader that returns uname data.
type Uname struct {
	Data map[string]interface{}
}

// Run gathers uname information from shell.
func (u *Uname) Run() error {
	cliBytes, err := exec.Command("uname", "-a").Output()
	if err != nil {
		return err
	}
	u.Data["Shell"] = strings.TrimSpace(string(cliBytes))

	return nil
}

// ToJson serialize Data field to JSON.
func (u *Uname) ToJson() ([]byte, error) {
	return json.Marshal(u.Data)
}
