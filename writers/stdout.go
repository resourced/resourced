package writers

import (
	"fmt"
)

// NewStdOut is StdOut constructor.
func NewStdOut() *StdOut {
	n := &StdOut{}
	return n
}

// StdOut is a writer that simply serialize all readers data to STDOUT.
type StdOut struct {
	Base
}

func (s *StdOut) Run() error {
	s.Data = s.GetReadersData()
	inJson, err := s.ToJson()

	if err != nil {
		return err
	}

	fmt.Println(string(inJson))
	return nil
}
