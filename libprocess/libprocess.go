package libprocess

import (
	"encoding/json"
	"github.com/chillaxio/chillax/libtime"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

func CheckProcessByPid(pid int) error {
	process, err := os.FindProcess(int(pid))
	if err == nil {
		err = process.Signal(syscall.Signal(0))
	}
	return err
}

func NewProcessWrapper(command string) *ProcessWrapper {
	p := &ProcessWrapper{
		Name:    path.Base(command),
		Command: command,
	}
	p.SetDefaults()
	return p
}

type ProcessWrapper struct {
	Name           string
	Command        string
	StopDelay      string
	StartDelay     string
	Ping           string
	Pid            int
	Status         string
	Respawn        int
	RespawnCounter int
	cmdStruct      *exec.Cmd
}

func (p *ProcessWrapper) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProcessWrapper) SetDefaults() {
	p.Ping = "30s"
	p.StopDelay = "0s"
	p.StartDelay = "0s"
	p.Pid = -1
	p.Respawn = -1
}

func (p *ProcessWrapper) NewCmd() *exec.Cmd {
	wd, _ := os.Getwd()

	parts := strings.Fields(p.Command)
	head := parts[0]
	parts = parts[1:len(parts)]

	cmd := exec.Command(head, parts...)
	cmd.Dir = wd
	cmd.Env = os.Environ()

	return cmd
}

func (p *ProcessWrapper) IsProcessStarted() bool {
	return p.cmdStruct.Process != nil
}

func (p *ProcessWrapper) IsProcessExists() bool {
	err := CheckProcessByPid(p.Pid)
	if err != nil && err.Error() == "no such process" {
		return false
	}
	return true
}

func (p *ProcessWrapper) StartAndWatch() error {
	err := p.Start()
	if err != nil {
		return err
	}

	p.DoPing(func() {
		if p.Pid > 0 {
			p.RespawnCounter = 0
			p.Status = "running"
		}
	})

	go p.Watch()

	return nil
}

// Start process
func (p *ProcessWrapper) Start() error {
	p.cmdStruct = p.NewCmd()

	err := p.cmdStruct.Start()
	if err != nil {
		return err
	}

	p.Pid = p.cmdStruct.Process.Pid
	p.Status = "started"

	p.ListenStopSignals()

	err = libtime.SleepString(p.StartDelay)

	return err
}

// Stop process and all its children
func (p *ProcessWrapper) Stop() error {
	var err error

	if p.cmdStruct != nil && p.cmdStruct.Process != nil {
		if p.Pid > 0 {
			err = p.cmdStruct.Process.Kill()
		}

		if err == nil {
			p.Release("stopped")
		}

		err = libtime.SleepString(p.StopDelay)
	}

	return err
}

// Release and remove process pidfile
func (p *ProcessWrapper) Release(status string) {
	if p.cmdStruct != nil && p.cmdStruct.Process != nil {
		p.cmdStruct.Process.Release()
	}
	p.Pid = -1
	p.Status = status
}

func (p *ProcessWrapper) RestartAndWatch() error {
	err := p.Stop()
	if err != nil {
		return err
	}

	err = p.StartAndWatch()
	if err != nil {
		return err
	}

	p.Status = "restarted"

	return nil
}

// Restart process
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (p *ProcessWrapper) Restart() error {
	err := p.Stop()
	if err != nil {
		return err
	}

	err = p.Start()
	if err != nil {
		return err
	}

	p.Status = "restarted"
	return nil
}

//Run callback on the process after *ProcessWrapper.Ping duration.
func (p *ProcessWrapper) DoPing(callback func()) {
	t, err := time.ParseDuration(p.Ping)
	if err == nil {
		go func() {
			select {
			case <-time.After(t):
				callback()
			}
		}()
	}
}

// Watch the process changes and restart if necessary
func (p *ProcessWrapper) Watch() {
	if p.cmdStruct.Process == nil {
		p.Release("stopped")
		return
	}

	procStateChan := make(chan *os.ProcessState)
	diedChan := make(chan error)

	go func() {
		state, err := p.cmdStruct.Process.Wait()
		if err != nil {
			diedChan <- err
			return
		}
		procStateChan <- state
	}()

	select {
	case <-procStateChan:
		if p.Status == "stopped" {
			return
		}

		p.RespawnCounter++

		if (p.Respawn != -1) && p.RespawnCounter > p.Respawn {
			p.Release("exited")
			return
		}

		p.RestartAndWatch()

	case <-diedChan:
		p.Release("killed")
	}
}

func (p *ProcessWrapper) ListenStopSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigChan

		if p.IsProcessStarted() {
			p.Stop()
			close(sigChan)
		}
	}()
}
