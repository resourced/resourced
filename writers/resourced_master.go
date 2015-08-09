package writers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
)

func init() {
	Register("ResourcedMasterHost", NewResourcedMasterHost)
	Register("ResourcedMasterStacks", NewResourcedMasterStacks)
}

// NewResourcedMasterHost is ResourcedMasterHost constructor.
func NewResourcedMasterHost() IWriter {
	return &ResourcedMasterHost{}
}

// ResourcedMasterHost is a writer that serialize readers data to ResourcedMasterHost.
type ResourcedMasterHost struct {
	Http
}

// NewResourcedMasterStacks is ResourcedMasterStacks constructor.
func NewResourcedMasterStacks() IWriter {
	return &ResourcedMasterStacks{}
}

type Stack struct {
	Steps []string `toml:"steps"`
}

// ResourcedMasterStacks is a writer that serialize readers data to ResourcedMasterStacks.
type ResourcedMasterStacks struct {
	Http
	Root       string
	CurrentSHA string
}

// stacksData gathers complete list of ResourceD Stacks metadata.
func (rm *ResourcedMasterStacks) stacksData() map[string]interface{} {
	data := make(map[string]interface{})
	logicList := make([]string, 0)
	stackList := make([]Stack, 0)

	if _, err := os.Stat(path.Join(rm.Root, "logic")); err == nil {
		logic, err := ioutil.ReadDir(path.Join(rm.Root, "logic"))
		if err == nil {
			for _, lgc := range logic {
				logicList = append(logicList, lgc.Name())
			}
		}
	}

	if _, err := os.Stat(path.Join(rm.Root, "stacks")); err == nil {
		stacks, err := ioutil.ReadDir(path.Join(rm.Root, "stacks"))
		if err == nil {
			for _, stackName := range stacks {
				var stk Stack

				stackPath := path.Join(rm.Root, "stacks", stackName.Name(), stackName.Name()+".toml")
				if _, err := toml.DecodeFile(stackPath, &stk); err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Errorf("Unable to decode %v", stackPath)
				}

				stackList = append(stackList, stk)
			}
		}
	}

	data["logic"] = logicList
	data["stacks"] = stackList

	return data
}

func (rm *ResourcedMasterStacks) GetCurrentSHA() string {
	output, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(string(output))

		return ""
	}
	return strings.TrimSpace(string(output))
}

// Run executes the writer.
func (rm *ResourcedMasterStacks) Run() error {
	if rm.Root == "" {
		return errors.New("ResourceD Stacks root should not be empty")
	}

	callback := func() error {
		rm.CurrentSHA = rm.GetCurrentSHA()

		data := rm.stacksData()
		data["CurrentSHA"] = rm.CurrentSHA

		dataJson, err := json.Marshal(data)
		if err != nil {
			return err
		}

		req, err := rm.NewHttpRequest(dataJson)
		if err != nil {
			return err
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error":      err.Error(),
				"req.URL":    req.URL.String(),
				"req.Method": req.Method,
			}).Error("Failed to send HTTP request")

			return err
		}

		if resp.Body != nil {
			resp.Body.Close()
		}

		// Set rm.Data
		rm.Data = data

		return nil
	}

	if rm.CurrentSHA == "" {
		return callback()
	}

	return rm.WatchDir(rm.Root, callback)
}
