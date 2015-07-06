package queryparser

import (
	"testing"
)

func TestJSONQuery(t *testing.T) {
	query := []byte(`[true]`)
	qp := New(query)

	jsonQuery, err := qp.JSONQuery()
	if err != nil {
		t.Fatalf("Unable to parse query. Error: %v", err)
	}

	if len(jsonQuery) != 1 {
		t.Errorf("Failed to parse query correctly. Length: %v", len(jsonQuery))
	}
}

func TestEvalSingleValue(t *testing.T) {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	query := []byte(`[true]`)
	qp := New(query)

	evaluated, err := qp.EvalExpressions(data, nil)
	if err != nil {
		t.Fatalf("Unable to evaluate query. Error: %v", err)
	}

	if evaluated != true {
		t.Errorf("Failed to parse query correctly. Value: %v", evaluated)
	}
}

func TestEvalSingleExpression(t *testing.T) {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	query := []byte(`[">", {"/r/load-avg": "LoadAvg1m"}, 0.5]`)
	qp := New(query)

	evaluated, err := qp.EvalExpressions(data, nil)
	if err != nil {
		t.Fatalf("Unable to evaluate query. Error: %v", err)
	}

	if evaluated != true {
		t.Errorf("Failed to parse query correctly. Value: %v", evaluated)
	}

}

func TestEvalBooleanExpressions(t *testing.T) {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	query := []byte(`["&&", [">", {"/r/load-avg": "LoadAvg1m"}, 0.5], ["<", {"/r/load-avg": "LoadAvg1m"}, 10]]`)
	qp := New(query)

	evaluated, err := qp.EvalExpressions(data, nil)
	if err != nil {
		t.Fatalf("Unable to evaluate query. Error: %v", err)
	}

	if evaluated != true {
		t.Errorf("Failed to parse query correctly. Value: %v", evaluated)
	}
}

func TestEvalNestedBooleanExpressions(t *testing.T) {
	data := make(map[string][]byte)
	data["/r/load-avg"] = []byte(`{"Data": {"LoadAvg1m": 0.904296875}}`)

	query := []byte(`["&&", ["&&", [">", {"/r/load-avg": "LoadAvg1m"}, 0.5], ["<", {"/r/load-avg": "LoadAvg1m"}, 10]], true]`)
	qp := New(query)

	evaluated, err := qp.EvalExpressions(data, nil)
	if err != nil {
		t.Fatalf("Unable to evaluate query. Error: %v", err)
	}

	if evaluated != true {
		t.Errorf("Failed to parse query correctly. Value: %v", evaluated)
	}
}
