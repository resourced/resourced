package writers

import (
	resourced_config "github.com/resourced/resourced/config"
	"os"
	"reflect"
	"testing"
)

func NewGoStructForTest(t *testing.T) IWriter {
	writer, err := NewGoStruct("StdOut")
	if err != nil {
		t.Fatalf("Creating a writer GoStruct should never fail. Error: %v", err)
	}
	return writer
}

func NewGoStructByConfigForTest(t *testing.T) IWriter {
	config, err := resourced_config.NewConfig(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/data/resourced-configs/writers/stdout.toml"), "writer")
	if err != nil {
		t.Fatalf("Initializing Config should work. Error: %v", err)
	}

	writer, err := NewGoStructByConfig(config)
	if err != nil {
		t.Fatalf("Initializing Writer should not fail. Error: %v", err)
	}

	for field, value := range map[string]string{
		"JsonProcessor": "$GOPATH/src/github.com/resourced/resourced/tests/data/script-writer/json-flattener.py"} {

		goStructField := reflect.ValueOf(writer).Elem().FieldByName(field)

		if goStructField.String() != value {
			t.Errorf("writer.%s is not set through the config. Value: %v, Writer: %v", field, goStructField.String(), writer)
		}
	}

	return writer
}

func TestSetReadersDataInBytes(t *testing.T) {
	writer := NewGoStructForTest(t)

	jsonData := `{
    "Data": {
        "LoadAvg15m": 1.59375,
        "LoadAvg1m": 1.5537109375,
        "LoadAvg5m": 1.68798828125
    },
    "GoStruct": "LoadAvg",
    "Host": {
        "Name":"MacBook-Pro.local",
        "Tags":[]
    },
    "Interval": "1s",
    "Path": "/load-avg",
    "Tags": [ ],
    "UnixNano": 1420607791403576000
}`
	readersData := make(map[string][]byte)
	readersData["/load-avg"] = []byte(jsonData)

	writer.SetReadersDataInBytes(readersData)

	key := "/load-avg"
	_, ok := writer.GetReadersData()[key]
	if !ok {
		t.Errorf("Key does not exist. Key: %v, Data: %v", key, writer.GetReadersData())
	}
}

func TestSetReadersData(t *testing.T) {
	writer := NewGoStructForTest(t)

	readersData := make(map[string]interface{})
	readersData["/load-avg"] = 1.5

	writer.SetReadersData(readersData)

	key := "/load-avg"
	_, ok := writer.GetReadersData()[key]
	if !ok {
		t.Errorf("Key does not exist. Key: %v, Data: %v", key, writer.GetReadersData())
	}
}

func TestJsonProcessor(t *testing.T) {
	writer := NewGoStructByConfigForTest(t)

	processor := writer.GetJsonProcessor()
	if processor == "" {
		t.Error("processor should not be blank. Processor: %v", processor)
	}
}
