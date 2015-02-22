package mysql

import (
	"strings"
	"testing"
)

func TestNewMysqlInformationSchemaTables(t *testing.T) {
	m := NewMysqlInformationSchemaTables()
	if m.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestMysqlInformationSchemaTablesRun(t *testing.T) {
	m := NewMysqlInformationSchemaTables()
	err := m.Run()
	if err != nil {
		t.Errorf("Fetching information_schema data should always be successful. Error: %v", err)
	}

	if len(m.Data["Tables"]) == 0 {
		jsonData, _ := m.ToJson()
		t.Errorf("Processlist data should never be empty. Data: %v", string(jsonData))
	}
}

func TestMysqlInformationSchemaTablesToJson(t *testing.T) {
	m := NewMysqlInformationSchemaTables()
	err := m.Run()
	if err != nil {
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
