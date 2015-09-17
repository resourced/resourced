package readers

import (
	"runtime"
	"testing"
)

func TestRunDarwin(t *testing.T) {
	if runtime.GOOS == "darwin" {
		err := NewIOStat().Run()
		if err == nil {
			t.Errorf("IOStat should only work on linux.")
		}
	}
}

func TestRunLinux(t *testing.T) {
	if runtime.GOOS == "linux" {
		reader := NewIOStat()

		err := reader.Run()
		if err != nil {
			t.Errorf("Reading iostat -x data should work on linux. Error: %v", err)
		}
	}
}
