// Package logline provides parsing logic on ResourceD TCP log line wire protocol.
// The basic protocol looks like this:
// type:base64|created:unix-timestamp|content:base64=
// type:plain|created:unix-timestamp|content:plaintext
package logline

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// Parse multiple lines of ResourceD TCP log line wire protocol.
func Parse(stringLoglines string) []LiveLogline {
	loglines := strings.Split(stringLoglines, "\n")
	result := make([]LiveLogline, len(loglines))

	for i, logline := range loglines {
		result[i] = ParseSingle(logline)
	}

	return result
}

// ParseSingle parses a single line of ResourceD TCP log line wire protocol.
func ParseSingle(logline string) LiveLogline {
	l := LiveLogline{}

	chunks := strings.Split(logline, "|")
	for _, chunk := range chunks {
		kv_slice := strings.Split(chunk, ":")

		if len(kv_slice) >= 2 {
			key := strings.TrimSpace(kv_slice[0])
			value := strings.TrimSpace(kv_slice[1])

			switch key {
			case "type":
				l.Type = value
			case "created":
				created, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					l.Created = created
				}
			case "content":
				l.Type = value
			}
		}
	}

	return l
}

type LiveLogline struct {
	Type    string
	Created int64
	Content string
}

// PlainContent returns the plaintext version of content.
func (l LiveLogline) PlainContent() string {
	if l.Type == "plain" {
		return l.Content
	}

	plain, err := base64.StdEncoding.DecodeString(l.Content)
	if err != nil {
		return "Failed to decode base64 content. Error: " + err.Error()
	}

	return plain
}

// Base64Content returns the base64 version of content.
func (l LiveLogline) Base64Content() string {
	if l.Type == "base64" {
		return l.Content
	}

	encoded, err := base64.StdEncoding.EncodeToString([]byte(l.Content))
	if err != nil {
		return "Failed to encode base64 content. Error: " + err.Error()
	}

	return encoded
}

// EncodePlain builds the wire protocol for plaintext type.
func (l LiveLogline) EncodePlain() string {
	return fmt.Sprintf("type:plain|created:%v|content:%v", l.Created, l.PlainContent())
}

// EncodePlain builds the wire protocol for base64 type.
func (l LiveLogline) EncodeBase64() string {
	return fmt.Sprintf("type:base64|created:%v|content:%v", l.Created, l.Base64Content())
}
