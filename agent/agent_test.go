package agent

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestConstructor(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("RESOURCED_TAGS", "prod, mysql, percona")
	os.Setenv("RESOURCED_CONFIG_READER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-reader")
	os.Setenv("RESOURCED_CONFIG_WRITER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-writer")

	agent, err := NewAgent()
	defer agent.Db.Close()

	if err != nil {
		t.Fatalf("Initializing ConfigStorage should work. Error: %v", err)
	}

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
	gopath := os.Getenv("GOPATH")
	os.Setenv("RESOURCED_CONFIG_READER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-reader")
	os.Setenv("RESOURCED_CONFIG_WRITER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-writer")

	agent, err := NewAgent()
	defer agent.Db.Close()

	if err != nil {
		t.Fatalf("Initializing ConfigStorage should work. Error: %v", err)
	}

	_, err = agent.Run(agent.ConfigStorage.Readers[1])
	if err != nil {
		t.Fatalf("Run should work. Error: %v", err)
	}
}

func TestGetRun(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("RESOURCED_CONFIG_READER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-reader")
	os.Setenv("RESOURCED_CONFIG_WRITER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-writer")

	agent, err := NewAgent()
	defer agent.Db.Close()

	if err != nil {
		t.Fatalf("Initializing ConfigStorage should work. Error: %v", err)
	}

	config := agent.ConfigStorage.Readers[1]

	_, err = agent.Run(config)
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
	gopath := os.Getenv("GOPATH")
	os.Setenv("RESOURCED_CONFIG_READER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-reader")
	os.Setenv("RESOURCED_CONFIG_WRITER_DIR", gopath+"/src/github.com/resourced/resourced/tests/data/config-writer")

	agent, _ := NewAgent()
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
