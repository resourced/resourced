package queryparser

import (
	"strings"
	"testing"
)

var queries = []string{
	`(((/r/load-avg.LoadAvg1m > 5) && (/r/load-avg.LoadAvg1m < 10)) || (/r/load-avg.LoadAvg1m == 100))`,
	`((/r/load-avg.LoadAvg1m > 5 && /r/load-avg.LoadAvg1m < 10) || (/r/load-avg.LoadAvg1m == 100))`,
	`(/r/load-avg.LoadAvg1m > 5 && /r/load-avg.LoadAvg1m < 10) || (/r/load-avg.LoadAvg1m == 100)`,
	`(/r/load-avg.LoadAvg1m > 5 && /r/load-avg.LoadAvg1m < 10) || /r/load-avg.LoadAvg1m == 100`,
	`/r/load-avg.LoadAvg1m > 5 && /r/load-avg.LoadAvg1m < 10 || /r/load-avg.LoadAvg1m == 100`,
	`/r/load-avg.LoadAvg1m > 5 && /r/load-avg.LoadAvg1m < 10 || false`,
	`/r/load-avg.LoadAvg1m > 5 && (true) || false`,
}

func TestDataValue(t *testing.T) {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	qp := New(data)

	valueInterface, err := qp.dataValue("/r/load-avg", "LoadAvg1m")
	if err != nil {
		t.Fatalf("Unable to fetch data value. Error: %v", err)
	}
	if valueInterface == nil {
		t.Fatalf("Data value should not be nil.")
	}
	if valueInterface.(float64) != 0.904296875 {
		t.Fatalf("Fetch data value incorrectly. Value: %v", valueInterface.(float64))
	}
}

func TestReplaceDataPathWithValue(t *testing.T) {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	qp := New(data)

	query, err := qp.replaceDataPathWithValue(queries[0])
	if err != nil {
		t.Fatalf("Failed to replace data path with value. Error: %v", err)
	}

	if strings.Contains(query, "/r/") {
		t.Fatalf("Failed to replace data path with value. Query: %v", query)
	}

}

func TestParse(t *testing.T) {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	qp := New(data)

	for _, query := range queries {
		result, err := qp.Parse(query)
		if err != nil {
			t.Fatalf("Failed to parse query. Error: %v", err)
		}

		if result != false {
			t.Fatalf("Failed to parse query correctly. Result: %v", result)
		}
	}
}
