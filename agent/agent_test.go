package agent

import (
	resourced_config "github.com/resourced/resourced/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func createAgentForAgentTest(t *testing.T) *Agent {
	os.Setenv("RESOURCED_CONFIG_READER_DIR", "$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader")
	os.Setenv("RESOURCED_CONFIG_WRITER_DIR", "$GOPATH/src/github.com/resourced/resourced/tests/data/config-writer")

	agent, err := NewAgent()
	if err != nil {
		t.Fatalf("Initializing agent should work. Error: %v", err)
	}
	return agent
}

func TestConstructor(t *testing.T) {
	os.Setenv("RESOURCED_TAGS", "prod, mysql, percona")

	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	if agent.DbPath == "" {
		t.Errorf("Default DbPath is set incorrectly. agent.DbPath: %v", agent.DbPath)
	}

	if _, err := os.Stat(agent.DbPath); err != nil {
		if os.IsNotExist(err) {
			t.Error("resourced directory does not exist.")
		}
	}

	if len(agent.Tags) != 3 {
		t.Error("agent.Tags should not be empty.")
	}

	for _, tag := range agent.Tags {
		if tag != "prod" && tag != "mysql" && tag != "percona" {
			t.Errorf("agent.Tags is incorrect: %v", agent.Tags)
		}
	}
}

func TestRun(t *testing.T) {
	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	_, err := agent.Run(agent.ConfigStorage.Readers[1])
	if err != nil {
		t.Fatalf("Run should work. Error: %v", err)
	}
}

func TestGetRun(t *testing.T) {
	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	config := agent.ConfigStorage.Readers[1]

	_, err := agent.Run(config)
	if err != nil {
		t.Fatalf("Run should work. Error: %v", err)
	}

	output, err := agent.GetRun(config)
	if err != nil {
		t.Fatalf("GetRun should work. Error: %v", err)
	}
	if string(output) == "" {
		t.Errorf("GetRun should return JSON data. Output: %v", string(output))
	}
}

func TestHttpRouter(t *testing.T) {
	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	_, err := agent.Run(agent.ConfigStorage.Readers[1])
	if err != nil {
		t.Fatalf("Run should work. Error: %v", err)
	}

	router := agent.HttpRouter()

	req, err := http.NewRequest("GET", "/r/cpu/info", nil)
	if err != nil {
		t.Errorf("Failed to create HTTP request. Error: %v", err)
	}

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if jsonData, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("Failed to read response body. Error: %v", err)
	} else {
		if strings.Contains(string(jsonData), "Error") {
			t.Errorf("jsonData shouldn't return error: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `UnixNano`) {
			t.Errorf("jsonData does not contain 'UnixNano' key: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Command`) && !strings.Contains(string(jsonData), `GoStruct`) {
			t.Errorf("jsonData does not contain 'Command' and 'GoStruct' keys: %s", jsonData)
		} else if !strings.Contains(string(jsonData), `Data`) {
			t.Errorf("jsonData does not contain 'Data' key: %s", jsonData)
		}
	}
}

func TestPathWithPrefix(t *testing.T) {
	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	config := agent.ConfigStorage.Readers[1]

	path := agent.pathWithPrefix(config)
	if !strings.HasPrefix(path, "/r") {
		t.Errorf("Path should have been prefixed with /r. Path: %v", path)
	}
	if strings.HasPrefix(path, "/w") {
		t.Errorf("Path is prefixed incorrectly. Path: %v", path)
	}
}

func TestPathWithReaderPrefix(t *testing.T) {
	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	toBeTested := agent.pathWithReaderPrefix("/stuff")
	if toBeTested != "/r/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = agent.pathWithReaderPrefix("/r/stuff")
	if toBeTested != "/r/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}
}

func TestInitGoStructReader(t *testing.T) {
	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	var config resourced_config.Config
	for _, c := range agent.ConfigStorage.Readers {
		if c.GoStruct == "DockerContainersMemory" {
			config = c
			break
		}
	}

	reader, err := agent.initGoStructReader(config)
	if err != nil {
		t.Fatalf("Initializing Reader should not fail. Error: %v", err)
	}

	goStructField := reflect.ValueOf(reader).Elem().FieldByName("CgroupBasePath")
	if goStructField.String() != "/sys/fs/cgroup/cpuacct/docker" {
		t.Errorf("reader.CgroupBasePath is not set through the config. CgroupBasePath: %v", goStructField.String())
	}
}

func TestInitGoStructWriter(t *testing.T) {
	agent := createAgentForAgentTest(t)
	defer agent.Db.Close()

	var config resourced_config.Config
	for _, c := range agent.ConfigStorage.Writers {
		if c.GoStruct == "Http" {
			config = c
			break
		}
	}

	writer, err := agent.initGoStructWriter(config)
	if err != nil {
		t.Fatalf("Initializing Writer should not fail. Error: %v", err)
	}

	for field, value := range map[string]string{"Url": "http://localhost:55556", "Method": "POST"} {
		goStructField := reflect.ValueOf(writer).Elem().FieldByName(field)
		if goStructField.String() != value {
			t.Errorf("writer.%s is not set through the config. Url: %v", field, goStructField.String())
		}
	}
}
