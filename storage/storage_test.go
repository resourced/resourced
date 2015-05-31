package storage

import (
	"testing"
)

func TestCrud(t *testing.T) {
	testData := `{"Free": 1000, "Used": 500}`

	s := NewStorage()
	s.Set("/free", []byte(testData))

	if string(s.Get("/free")) != testData {
		t.Errorf("Failed to test set and get. Actual Data: %v", string(s.Get("/free")))
	}
}
