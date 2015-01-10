package writers

// NewNoop is Noop constructor.
func NewNoop() *Noop {
	n := &Noop{}
	return n
}

// Noop is a writer that does not do anything.
type Noop struct {
	Base
}
