package mysql

import (
	"strings"
	"testing"
)

func TestMysqlProcesslistRun(t *testing.T) {
	m := &MysqlProcesslist{}
	m.Data = make(map[string][]Processlist)
	m.Retries = 1

	err := m.Run()
	if strings.Contains(err.Error(), "connection refused") {
		return
	}

	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("Fetching processlist data should always be successful. Error: %v", err)
	}

	if len(m.Data) == 0 {
		jsonData, _ := m.ToJson()
		t.Errorf("Processlist data should never be empty. Data: %v", string(jsonData))
	}
}

func TestMysqlProcesslistToJson(t *testing.T) {
	m := &MysqlProcesslist{}
	m.Data = make(map[string][]Processlist)
	m.Retries = 1

	err := m.Run()
	if strings.Contains(err.Error(), "connection refused") {
		t.Infof("Local MySQL is not running. Stop testing.")
		return
	}

	if err != nil && !strings.Contains(err.Error(), "connection refused") {
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
