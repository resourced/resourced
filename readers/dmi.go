package readers

import (
	"encoding/json"

	"github.com/dselans/dmidecode"
)

func init() {
	Register("DMI", NewDMI)
}

func NewDMI() IReader {
	d := &DMI{}
	d.Data = make(map[string]map[string]string)
	return d
}

type DMI struct {
	Data map[string]map[string]string
}

func (d *DMI) Run() error {
	dmi := dmidecode.New()

	if err := dmi.Run(); err != nil {
		return err
	}

	d.Data = dmi.Data

	return nil
}

func (d *DMI) ToJson() ([]byte, error) {
	return json.Marshal(d.Data)
}
