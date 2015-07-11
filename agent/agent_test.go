package agent

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	resourced_config "github.com/resourced/resourced/config"
)

func createAgentForTest(t *testing.T) *Agent {
	os.Setenv("RESOURCED_CONFIG_DIR", os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/data/resourced-configs"))

	// Provide empty slice - allow all to connect
	agent, err := NewAgent([]*net.IPNet{})
	if err != nil {
		t.Fatalf("Initializing agent should work. Error: %v", err)
	}

	return agent
}

func TestConstructor(t *testing.T) {
	os.Setenv("RESOURCED_TAGS", "prod, mysql, percona")

	agent := createAgentForTest(t)

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
	agent := createAgentForTest(t)

	_, err := agent.Run(agent.ConfigStorage.Readers[1])
	if err != nil {
		t.Fatalf("Run should work. Error: %v", err)
	}
}

func TestGetRun(t *testing.T) {
	agent := createAgentForTest(t)

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
	agent := createAgentForTest(t)

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
			t.Errorf("jsonData shouldn't return error: %s, %s", jsonData, req.RemoteAddr)
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
	agent := createAgentForTest(t)

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
	agent := createAgentForTest(t)

	toBeTested := agent.pathWithReaderPrefix("/stuff")
	if toBeTested != "/r/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = agent.pathWithReaderPrefix("/r/stuff")
	if toBeTested != "/r/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}
}

func TestPathWithWriterPrefix(t *testing.T) {
	agent := createAgentForTest(t)

	toBeTested := agent.pathWithWriterPrefix("/stuff")
	if toBeTested != "/w/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}

	toBeTested = agent.pathWithWriterPrefix("/w/stuff")
	if toBeTested != "/w/stuff" {
		t.Errorf("Path is prefixed incorrectly. toBeTested: %v", toBeTested)
	}
}

func TestInitGoStructReader(t *testing.T) {
	agent := createAgentForTest(t)

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
	if goStructField.String() != "/sys/fs/cgroup/memory/docker" {
		t.Errorf("reader.CgroupBasePath is not set through the config. CgroupBasePath: %v", goStructField.String())
	}
}

func TestInitGoStructWriter(t *testing.T) {
	agent := createAgentForTest(t)

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

	for field, value := range map[string]string{
		"Url":     "http://localhost:55655/",
		"Method":  "POST",
		"Headers": "X-Token=abc123,X-Teapot-Count=2"} {

		goStructField := reflect.ValueOf(writer).Elem().FieldByName(field)
		if goStructField.String() != value {
			t.Errorf("writer.%s is not set through the config. Url: %v", field, goStructField.String())
		}
	}
}

func TestCommonData(t *testing.T) {
	agent := createAgentForTest(t)

	var config resourced_config.Config
	for _, c := range agent.ConfigStorage.Readers {
		if c.GoStruct == "DockerContainersMemory" {
			config = c
			break
		}
	}

	record := agent.commonData(config)
	if len(record) == 0 {
		t.Error("common data should never be empty")
	}
	for _, key := range []string{"UnixNano", "Path", "Interval"} {
		if _, ok := record[key]; !ok {
			t.Errorf("%v data should never be empty.", key)
		}
	}
}

func TestIsAllowed(t *testing.T) {
	_, network, _ := net.ParseCIDR("127.0.0.1/8")
	allowedNetworks := []*net.IPNet{network}

	agent, err := NewAgent(allowedNetworks)
	if err != nil {
		t.Fatalf("Initializing agent should work. Error: %v", err)
	}

	goodIP := "127.0.0.1"
	badIP := "10.0.0.1"
	brokenIP := "batman"

	if !agent.IsAllowed(goodIP) {
		t.Errorf("'%s' should be allowed", goodIP)
	}

	if agent.IsAllowed(badIP) {
		t.Errorf("'%s' should not be allowed", badIP)
	}

	if agent.IsAllowed(brokenIP) {
		t.Errorf("Invalid IP address '%s' should not be allowed ", brokenIP)
	}
}
