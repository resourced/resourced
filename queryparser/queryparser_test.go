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

var tagQueries = []string{
	`tags.role == "appserver" && tags.environment == "prod"`,
	`tags.role == "appserver" && tags.environment == "staging"`,
	`(tags.role == "appserver" && tags.environment == "staging")`,
}

var nameQueries = []string{
	`name == "some-hostname"`,
}

var metadataQueries = []string{
	`metadata.users/didip.name == "didip"`,
	`metadata.users/didip.uid == 1000`,
}

func queryparserForTest(t *testing.T) *QueryParser {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	tags := make(map[string]string)
	tags["role"] = "appserver"
	tags["environment"] = "staging"

	qp := New(data, tags)

	return qp
}

func TestDataValue(t *testing.T) {
	qp := queryparserForTest(t)

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
	qp := queryparserForTest(t)

	query, err := qp.replaceDataPathWithValue(queries[0])
	if err != nil {
		t.Fatalf("Failed to replace data path with value. Error: %v", err)
	}

	if strings.Contains(query, "/r/") {
		t.Fatalf("Failed to replace data path with value. Query: %v", query)
	}
}

func TestReplaceTagsWithValue(t *testing.T) {
	qp := queryparserForTest(t)

	query, err := qp.replaceTagsWithValue(tagQueries[0])
	if err != nil {
		t.Fatalf("Failed to replace tags with value. Error: %v", err)
	}

	if strings.Contains(query, "tags.") {
		t.Fatalf("Failed to replace tags with value. Query: %v", query)
	}
}

func TestReplaceHostnameWithValue(t *testing.T) {
	qp := queryparserForTest(t)

	query, err := qp.replaceHostnameWithValue(nameQueries[0])
	if err != nil {
		t.Fatalf("Failed to replace tags with value. Error: %v", err)
	}
	if strings.Contains(query, "name ==") {
		t.Fatalf("Failed to replace tags with value. Query: %v", query)
	}
}

func TestParseQueries(t *testing.T) {
	qp := queryparserForTest(t)

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

func TestParseTagQueries(t *testing.T) {
	qp := queryparserForTest(t)

	for i, query := range tagQueries {
		result, err := qp.Parse(query)
		if err != nil {
			t.Fatalf("Failed to parse tag query. Error: %v", err)
		}

		if i == 0 {
			if result != false {
				t.Fatalf("Failed to parse tag query correctly. Result: %v", result)
			}

		} else {
			if result != true {
				t.Fatalf("Failed to parse tag query correctly. Result: %v", result)
			}
		}
	}
}

func TestParseHostnameQueries(t *testing.T) {
	qp := queryparserForTest(t)

	for i, query := range nameQueries {
		result, err := qp.Parse(query)
		if err != nil {
			t.Fatalf("Failed to parse tag query. Error: %v", err)
		}

		if i == 0 {
			if result != false {
				t.Fatalf("Failed to parse tag query correctly. Result: %v", result)
			}
		}
	}
}
