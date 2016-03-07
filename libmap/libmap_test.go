package libmap

import (
	"testing"
)

func TestTSafeMapBytesSetGet(t *testing.T) {
	testData := `{"Free": 1000, "Used": 500}`

	s := NewTSafeMapBytes()
	s.Set("/free", []byte(testData))

	if string(s.Get("/free")) != testData {
		t.Errorf("Failed to test set and get. Actual Data: %v", string(s.Get("/free")))
	}
}
