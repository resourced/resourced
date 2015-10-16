package libstring

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestReplaceTildeWithRoot(t *testing.T) {
	path := "~/resourced"
	toBeTested := strings.Replace(path, "~", "/root", 1)

	if toBeTested != "/root/resourced" {
		t.Errorf("~ is not expanded correctly. Path: %v", toBeTested)
	}
}

func TestExpandTildeAndEnv(t *testing.T) {
	toBeTested := ExpandTildeAndEnv("~/resourced")

	if runtime.GOOS == "darwin" {
		if !strings.HasPrefix(toBeTested, "/Users") {
			t.Errorf("~ is not expanded correctly. Path: %v", toBeTested)
		}
	}

	toBeTested = ExpandTildeAndEnv("$GOPATH/src/github.com/resourced/resourced/tests/script-reader/darwin-memory.py")
	gopath := os.Getenv("GOPATH")

	if !strings.HasPrefix(toBeTested, gopath) {
		t.Errorf("$GOPATH is not expanded correctly. Path: %v", toBeTested)
	}
}

func TestGeneratePassword(t *testing.T) {
	_, err := GeneratePassword(8)
	if err != nil {
		t.Errorf("Generating password should not fail. err: %v", err)
	}
}

func TestGetIP(t *testing.T) {
	goodAddress := "127.0.0.1:55555"
	badAddress := "tasty:cakes"

	goodIP := GetIP(goodAddress)
	if goodIP == nil {
		t.Error("Should be able to parse '%v'", goodAddress)
	}

	if goodIP.String() != strings.Split(goodAddress, ":")[0] {
		t.Error("goodIP.String() should be the same as split goodAddress")
	}

	badIP := GetIP(badAddress)
	if badIP != nil {
		t.Error("Should not be able to parse '%v'", badAddress)
	}
}
