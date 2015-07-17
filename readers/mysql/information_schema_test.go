package mysql

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMysqlInformationSchemaTablesRun(t *testing.T) {
	m := &MysqlInformationSchemaTables{}
	m.Data = make(map[string][]InformationSchemaTables)
	m.Retries = 1

	err := m.Run()
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		t.Logf("Local MySQL is not running. Stop testing.")
		return
	}

	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("Fetching information_schema data should always be successful. Error: %v", err)
	}

	var data map[string][]InformationSchemaTables
	inJson, _ := m.ToJson()
	json.Unmarshal(inJson, &data)

	if len(data["Tables"]) == 0 {
		t.Errorf("Processlist data should never be empty. Data: %v", string(inJson))
	}
}

func TestMysqlInformationSchemaTablesToJson(t *testing.T) {
	m := NewMysqlInformationSchemaTables()
	err := m.Run()
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		t.Logf("Local MySQL is not running. Stop testing.")
		return
	}

	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("Fetching information_schema data should always be successful. Error: %v", err)
	}

	jsonData, err := m.ToJson()
	if err != nil {
		t.Errorf("Marshalling information_schema data should always be successful. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}
}
