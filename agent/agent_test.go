package agent

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	resourced_config "github.com/resourced/resourced/config"
	_ "github.com/resourced/resourced/readers/docker"
)

func createAgentForTest(t *testing.T) *Agent {
	agent := createAgentWithAccessTokensForTest(t)

	agent.AccessTokens = make([]string, 0)

	return agent
}

func createAgentWithAccessTokensForTest(t *testing.T) *Agent {
	agent, err := New(os.ExpandEnv("$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs"))
	if err != nil {
		t.Fatalf("Initializing agent should work. Error: %v", err)
	}

	return agent
}

func TestRun(t *testing.T) {
	agent := createAgentForTest(t)

	if len(agent.Configs.Readers) <= 0 {
		t.Fatalf("Agent config readers should exist")
	}

	_, err := agent.Run(agent.Configs.Readers[1])
	if err != nil {
		t.Fatalf("Run should work. Error: %v", err)
	}
}

func TestHttpRouter(t *testing.T) {
	agent := createAgentForTest(t)

	_, err := agent.Run(agent.Configs.Readers[0])
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
		}
	}
}

func TestInitGoStructReader(t *testing.T) {
	agent := createAgentForTest(t)

	var config resourced_config.Config
	for _, c := range agent.Configs.Readers {
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
	for _, c := range agent.Configs.Writers {
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
