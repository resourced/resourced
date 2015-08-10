package executors

import (
	"encoding/json"
	"os"

	resourced_stacks_engine "github.com/resourced/resourced-stacks/engine"
)

func init() {
	Register("ResourcedStacks", NewResourcedStacks)
}

func NewResourcedStacks() IExecutor {
	s := &ResourcedStacks{}
	s.Data = make(map[string]interface{})

	return s
}

type ResourcedStacks struct {
	Base
	Name      string
	Root      string
	GitRepo   string
	GitBranch string
	DryRun    bool
	Data      map[string]interface{}
}

func (s *ResourcedStacks) onError(err error) error {
	s.Data["Error"] = err.Error()
	s.Data["ExitStatus"] = 1
	return err
}

// Run ResourceD Stacks and store the output on c.Data.
func (s *ResourcedStacks) Run() error {
	if s.IsConditionMet() && !s.HighThresholdExceeded() {
		s.Data["Conditions"] = s.Conditions
		s.Data["Name"] = s.Name
		s.Data["Root"] = s.Root
		s.Data["GitRepo"] = s.GitRepo
		s.Data["GitBranch"] = s.GitBranch

		// Create root directory
		err := os.MkdirAll(s.Root, 0755)
		if err != nil {
			return s.onError(err)
		}

		// Construct engine struct
		engine, err := resourced_stacks_engine.New(s.Root, s.Conditions)
		if err != nil {
			return s.onError(err)
		}

		engine.DryRun = s.DryRun
		engine.Git.HTTPS = s.GitRepo

		if s.GitBranch == "" {
			s.GitBranch = "master"
		}
		engine.Git.Branch = s.GitBranch

		// Perform git clone or pull
		err = engine.GitPull()
		if err != nil {
			return s.onError(err)
		}

		// Run stack
		output, err := engine.RunStack(s.Name, nil)

		s.Data["Output"] = string(output)

		if err != nil {
			return s.onError(err)
		} else {
			s.Data["ExitStatus"] = 0
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (s *ResourcedStacks) ToJson() ([]byte, error) {
	return json.Marshal(s.Data)
}
