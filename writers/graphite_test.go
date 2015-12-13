package writers

import (
	"encoding/json"
	"testing"
)

func newWriterForGraphiteTest() *Graphite {
	g := &Graphite{}
	g.Host = "localhost"
	g.Port = 2003
	return g
}

func TestToJson(t *testing.T) {
	g := newWriterForGraphiteTest()

	data := make(map[string]map[string]bool)
	data["parent"] = make(map[string]bool)
	data["parent"]["child1"] = true
	data["parent"]["child2"] = false

	g.Data = data

	inJson, err := g.ToJson()
	if err != nil {
		t.Errorf("Flattening data to JSON should not fail. Error: %v", err)
	}

	flattened := make(map[string]bool)

	err = json.Unmarshal(inJson, &flattened)
	if err != nil {
		t.Errorf("Flattening data to JSON should not fail. Error: %v", err)
	}
	if flattened["parent.child1"] != true {
		t.Errorf("Flattened data incorrectly. Result: %v", flattened["parent.child1"])
	}
	if flattened["parent.child2"] != false {
		t.Errorf("Flattened data incorrectly. Result: %v", flattened["parent.child2"])
	}
}
