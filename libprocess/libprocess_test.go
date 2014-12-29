package libprocess

import (
	"encoding/json"
	"github.com/chillaxio/chillax/libtime"
	"testing"
)

func NewProcessWrapperForTest() *ProcessWrapper {
	p := &ProcessWrapper{
		Name:    "SimpleHTTPServer",
		Command: "/usr/bin/python -m SimpleHTTPServer",
	}
	p.SetDefaults()
	p.Ping = "10ms"
	return p
}

func CheckBasicStartForTest(t *testing.T, p *ProcessWrapper) {
	err := p.Start()
	if err != nil {
		t.Errorf("Unable to start process. Error: %v", err)
	}
	if p.Status != "started" {
		t.Errorf("process status is set incorrectly. Process status: %v", p.Status)
	}
	if p.Pid <= 0 || p.cmdStruct.Process.Pid <= 0 {
		t.Errorf("Process should start with PID > 0. p.Pid: %v, p.cmdStruct.Process.Pid: %v", p.Pid, p.cmdStruct.Process.Pid)
	}
	if p.Pid != p.cmdStruct.Process.Pid {
		t.Errorf("ProcessWrapper PID should == Process PID")
	}
}

func CheckStartAndWatchForTest(t *testing.T, p *ProcessWrapper) {
	err := p.StartAndWatch()
	if err != nil {
		t.Errorf("Unable to start process. Error: %v", err)
	}
	if p.Status != "started" {
		t.Errorf("process status is set incorrectly. Process status: %v", p.Status)
	}
	if p.Pid <= 0 || p.cmdStruct.Process.Pid <= 0 {
		t.Errorf("Process should start with PID > 0. p.Pid: %v, p.cmdStruct.Process.Pid: %v", p.Pid, p.cmdStruct.Process.Pid)
	}
	if p.Pid != p.cmdStruct.Process.Pid {
		t.Errorf("ProcessWrapper PID should == Process PID")
	}
}

func CheckBasicStopForTest(t *testing.T, p *ProcessWrapper) {
	err := p.Stop()
	if err != nil {
		t.Errorf("Unable to stop process. Error: %v", err)
	}
	if p.Status != "stopped" {
		t.Errorf("process status is set incorrectly. Process status: %v", p.Status)
	}
}

func TestToJson(t *testing.T) {
	p := NewProcessWrapperForTest()

	CheckBasicStartForTest(t, p)

	inJson, _ := p.ToJson()

	var deserializedData map[string]interface{}

	err := json.Unmarshal(inJson, &deserializedData)
	if err != nil {
		t.Errorf("Unable to deserialize JSON. Error: %v", err)
	}

	if deserializedData["Name"].(string) != p.Name {
		t.Errorf("Bad deserialization")
	}

	CheckBasicStopForTest(t, p)
}

func TestProcessStartRestartStop(t *testing.T) {
	p := NewProcessWrapperForTest()

	CheckBasicStartForTest(t, p)

	err := p.Restart()
	if err != nil {
		t.Errorf("Unable to restart process. Error: %v", err)
	}
	if p.Status != "restarted" {
		t.Errorf("process status is set incorrectly. Process status: %v", p.Status)
	}

	CheckBasicStopForTest(t, p)
}

func TestProcessStartAndWatch(t *testing.T) {
	p := NewProcessWrapperForTest()

	CheckStartAndWatchForTest(t, p)

	libtime.SleepString("14ms")

	if p.Status != "running" {
		t.Errorf("process status is set incorrectly. Process status: %v", p.Status)
	}

	firstPid := p.Pid

	if p.cmdStruct.Process.Pid > 0 {
		err := p.cmdStruct.Process.Kill()
		if err != nil {
			t.Errorf("Unable to kill process manually.")
		}
		libtime.SleepString("120ms")

		secondPid := p.Pid

		if firstPid == secondPid {
			t.Errorf("New process should have generated new PID. firstPid: %v, secondPid: %v", firstPid, secondPid)
		}
	}

	CheckBasicStopForTest(t, p)
}

func TestProcessRestartAndWatch(t *testing.T) {
	p := NewProcessWrapperForTest()

	CheckStartAndWatchForTest(t, p)

	libtime.SleepString("14ms")

	if p.Status != "running" {
		t.Errorf("process status is set incorrectly. Process status: %v", p.Status)
	}

	firstPid := p.Pid

	err := p.RestartAndWatch()
	if err != nil {
		t.Errorf("Unable to restart and watch process. Error: %v", err)
	}

	if p.Status != "restarted" {
		t.Errorf("process status is set incorrectly. Process status: %v", p.Status)
	}

	libtime.SleepString("14ms")

	secondPid := p.Pid

	if firstPid == secondPid {
		t.Errorf("New process should have generated new PID. firstPid: %v, secondPid: %v", firstPid, secondPid)
	}

	CheckBasicStopForTest(t, p)
}
