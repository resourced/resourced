package mysql

import (
	"strings"
	"testing"
)

func TestNewMysqlProcesslist(t *testing.T) {
	m := NewMysqlProcesslist()
	if m.Data == nil {
		t.Error("Reader data should never be nil.")
	}
}

func TestMysqlProcesslistRun(t *testing.T) {
	m := NewMysqlProcesslist()
	err := m.Run()
	if err != nil {
		t.Errorf("Fetching processlist data should always be successful. Error: %v", err)
	}

	if len(m.Data["Processes"]) == 0 {
		jsonData, _ := m.ToJson()
		t.Errorf("Processlist data should never be empty. Data: %v", string(jsonData))
	}
}

func TestMysqlProcesslistToJson(t *testing.T) {
	m := NewMysqlProcesslist()
	err := m.Run()
	if err != nil {
		t.Errorf("Fetching processlist data should always be successful. Error: %v", err)
	}

	jsonData, err := m.ToJson()
	if err != nil {
		t.Errorf("Marshalling processlist data should always be successful. Error: %v", err)
	}

	jsonDataString := string(jsonData)

	if strings.Contains(jsonDataString, "Error") {
		t.Errorf("jsonDataString shouldn't return error: %v", jsonDataString)
	}
}
