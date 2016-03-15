package libmap

import (
	"testing"
)

func TestTSafeMapBytesSetGet(t *testing.T) {
	testData := `{"Free": 1000, "Used": 500}`

	s := NewTSafeMapBytes(nil)
	s.Set("/free", []byte(testData))

	if string(s.Get("/free")) != testData {
		t.Errorf("Failed to test set and get. Actual Data: %v", string(s.Get("/free")))
	}
}

func TestTSafeNestedMapInterfaceInitNestedMap(t *testing.T) {
	m := NewTSafeNestedMapInterface(nil)
	m.initNestedMap("aaa.bbb.ccc")

	if m.Get("aaa") == nil {
		t.Fatalf("Failed to init nested map")
	}
	if m.Get("aaa").(map[string]interface{})["bbb"] == nil {
		t.Fatalf("Failed to init nested map")
	}

	m.Get("aaa").(map[string]interface{})["bbb"].(map[string]interface{})["ccc"] = 42

	val := m.Get("aaa").(map[string]interface{})["bbb"].(map[string]interface{})["ccc"].(int)
	expected := 42
	if val != expected {
		t.Fatalf("Failed to get value on nested map. Expected: %v, Got: %v", expected, val)
	}
}

func TestTSafeNestedMapInterfaceSetGet(t *testing.T) {
	m := NewTSafeNestedMapInterface(nil)

	m.Set("aaa.bbb.ccc", 42)

	if m.Get("aaa") == nil {
		t.Fatalf("Failed to init nested map")
	}
	if m.Get("aaa").(map[string]interface{})["bbb"] == nil {
		t.Fatalf("Failed to init nested map")
	}

	val := m.Get("aaa").(map[string]interface{})["bbb"].(map[string]interface{})["ccc"].(int)
	expected := 42
	if val != expected {
		t.Fatalf("Failed to get value on nested map. Expected: %v, Got: %v", expected, val)
	}
}

func TestNewTSafeMapStringsBasicFunctionality(t *testing.T) {
	logDB := NewTSafeMapStrings(map[string][]string{
		"Loglines": make([]string, 0),
	})

	logDB.Append("Loglines", "some log")

	logs := logDB.Get("Loglines")
	if len(logs) != 1 {
		t.Fatalf("Failed to get value on string slice. Expected: %v, Got: %v", 1, len(logs))
	}

	logDB.Reset("Loglines")

	logs = logDB.Get("Loglines")
	if len(logs) != 0 {
		t.Fatalf("Failed to get value on string slice. Expected: %v, Got: %v", 0, len(logs))
	}
}
