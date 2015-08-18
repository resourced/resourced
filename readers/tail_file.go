package readers

import (
	"encoding/json"
	"errors"

	"github.com/ActiveState/tail"
	"github.com/resourced/resourced/libstring"
)

func init() {
	Register("TailFile", NewTailFile)
}

// NewTailFile is TailFile constructor.
func NewTailFile() IReader {
	c := &TailFile{}
	c.Data = make(map[string]interface{})
	return c
}

// TailFile is a reader that can watch a file.
type TailFile struct {
	FilePath string
	Data     map[string]interface{}
}

// Run executes the reader.
func (tl *TailFile) Run() error {
	return nil
}

// Tailer returns tail struct to watch file.
func (tl *TailFile) Tailer() (*tail.Tail, error) {
	if tl.FilePath == "" {
		return nil, errors.New("FilePath must not be empty")
	}

	tl.FilePath = libstring.ExpandTildeAndEnv(tl.FilePath)

	return tail.TailFile(tl.FilePath, tail.Config{
		Follow: true,
		ReOpen: true})
}

func (c *TailFile) ToJson() ([]byte, error) {
	return json.Marshal(c.Data)
}
