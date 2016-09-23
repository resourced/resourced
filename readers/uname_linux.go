// +build linux
package readers

import (
	"encoding/json"
	"os/exec"
	"syscall"
)

func init() {
	Register("Uname", NewUname)
}

func NewUname() IReader {
	u := &Uname{}
	u.Data = make(map[string]interface{})
	return u
}

// helper function to convert the unfortunate uname type.
func charsToString(ca *[65]int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}

// Uname is a reader that returns uname data.
type Uname struct {
	Data map[string]interface{}
}

// Run gathers uname information from syscall.
func (u *Uname) Run() error {
	var n syscall.Utsname

	err := syscall.Uname(&n)
	if err != nil {
		return err
	}

	u.Data["Sysname"] = charsToString(&n.Sysname)
	u.Data["Nodename"] = charsToString(&n.Nodename)
	u.Data["Release"] = charsToString(&n.Release)
	u.Data["Version"] = charsToString(&n.Version)
	u.Data["Machine"] = charsToString(&n.Machine)
	u.Data["Domainname"] = charsToString(&n.Domainname)

	cliBytes, err := exec.Command("uname", "-a").Output()
	if err != nil {
		return err
	}
	u.Data["Shell"] = string(cliBytes)

	return nil
}

// ToJson serialize Data field to JSON.
func (u *Uname) ToJson() ([]byte, error) {
	return json.Marshal(u.Data)
}
