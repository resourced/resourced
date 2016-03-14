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

	if m.Data["aaa"] == nil {
		t.Fatalf("Failed to init nested map")
	}
	if m.Data["aaa"].(map[string]interface{})["bbb"] == nil {
		t.Fatalf("Failed to init nested map")
	}

	m.Data["aaa"].(map[string]interface{})["bbb"].(map[string]interface{})["ccc"] = 42

	val := m.Data["aaa"].(map[string]interface{})["bbb"].(map[string]interface{})["ccc"].(int)
	expected := 42
	if val != expected {
		t.Fatalf("Failed to get value on nested map. Expected: %v, Got: %v", expected, val)
	}
}

func TestTSafeNestedMapInterfaceSetGet(t *testing.T) {
	m := NewTSafeNestedMapInterface(nil)

	m.Set("aaa.bbb.ccc", 42)

	if m.Data["aaa"] == nil {
		t.Fatalf("Failed to init nested map")
	}
	if m.Data["aaa"].(map[string]interface{})["bbb"] == nil {
		t.Fatalf("Failed to init nested map")
	}

	val := m.Data["aaa"].(map[string]interface{})["bbb"].(map[string]interface{})["ccc"].(int)
	expected := 42
	if val != expected {
		t.Fatalf("Failed to get value on nested map. Expected: %v, Got: %v", expected, val)
	}
}
