package logline

import (
	"strings"
	"testing"
)

func loglineForTest() string {
	return "type:plain|created:1468640833|content:Saturday, 16-Jul-16 03:47:13 UTC hello world log"
}

func TestParseSingle(t *testing.T) {
	lg := ParseSingle(loglineForTest())
	if lg.Type != "plain" {
		t.Errorf("Failed to assign type. Type: %v", lg.Type)
	}
	if lg.Created != int64(1468640833) {
		t.Errorf("Failed to assign created timestamp. Created: %v", lg.Created)
	}
	if lg.Content != "Saturday, 16-Jul-16 03:47:13 UTC hello world log" {
		t.Errorf("Failed to assign content. Content: %v", lg.Content)
	}
}

func TestParseMultiple(t *testing.T) {
	loglines := strings.Join([]string{loglineForTest(), loglineForTest(), loglineForTest()}, "\n")

	lgs := Parse(loglines)
	for _, lg := range lgs {
		if lg.Type != "plain" {
			t.Errorf("Failed to assign type. Type: %v", lg.Type)
		}
		if lg.Created != int64(1468640833) {
			t.Errorf("Failed to assign created timestamp. Created: %v", lg.Created)
		}
		if lg.Content != "Saturday, 16-Jul-16 03:47:13 UTC hello world log" {
			t.Errorf("Failed to assign content. Content: %v", lg.Content)
		}
	}
}

func TestDifferentContent(t *testing.T) {
	lg := ParseSingle(loglineForTest())
	if lg.PlainContent() != lg.Content {
		t.Errorf("Failed to return the correct content: %v", lg.PlainContent())
	}
	if lg.Base64Content() != "U2F0dXJkYXksIDE2LUp1bC0xNiAwMzo0NzoxMyBVVEMgaGVsbG8gd29ybGQgbG9n" {
		t.Errorf("Failed to return the correct content: %v", lg.Base64Content())
	}
}

func TestEncodePlain(t *testing.T) {
	lg := ParseSingle(loglineForTest())
	if lg.EncodePlain() != loglineForTest() {
		t.Errorf("Failed to encode correctly: %v", lg.EncodePlain())
	}
}

func TestEncodeBase64(t *testing.T) {
	lg := ParseSingle(loglineForTest())
	if lg.EncodeBase64() != "type:base64|created:1468640833|content:U2F0dXJkYXksIDE2LUp1bC0xNiAwMzo0NzoxMyBVVEMgaGVsbG8gd29ybGQgbG9n" {
		t.Errorf("Failed to encode correctly: %v", lg.EncodeBase64())
	}
}
