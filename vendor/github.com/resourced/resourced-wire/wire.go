// Package wire provides parsing logic on ResourceD TCP wire protocol.
// The protocol looks like this:
//     type:base64|created:unix-timestamp|content:base64=
//     type:plain|created:unix-timestamp|content:plaintext
//     type:json|created:unix-timestamp|content:{"foo": "bar"}
//
// Within Master daemon, this wire format is used for passing data through the message bus.
// The protocol looks like this:
//     topic:topic-name|type:json|created:unix-timestamp|content:{"foo": "bar"}
//
package wire

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// Parse multiple lines of ResourceD TCP log line wire protocol.
func Parse(stringLoglines string) []Wire {
	loglines := strings.Split(stringLoglines, "\n")
	result := make([]Wire, len(loglines))

	for i, logline := range loglines {
		result[i] = ParseSingle(logline)
	}

	return result
}

// ParseSingle parses a single line of ResourceD TCP log line wire protocol.
func ParseSingle(logline string) Wire {
	l := Wire{}

	chunks := strings.Split(logline, "|")
	for _, chunk := range chunks {
		kv_slice := strings.Split(chunk, ":")

		if len(kv_slice) >= 2 {
			key := strings.TrimSpace(kv_slice[0])
			value := strings.TrimSpace(kv_slice[1])

			switch key {
			case "topic":
				l.Topic = value
			case "type":
				l.Type = value
			case "created":
				created, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					l.Created = created
				}
			case "content":
				l.Content = strings.Join(kv_slice[1:], ":")
			}
		}
	}

	return l
}

type Wire struct {
	Topic   string
	Type    string
	Created int64
	Content string
}

// PlainContent returns the plaintext version of content.
func (l Wire) PlainContent() string {
	if l.Type == "plain" {
		return l.Content
	}

	plain, err := base64.StdEncoding.DecodeString(l.Content)
	if err != nil {
		return "Failed to decode base64 content. Error: " + err.Error()
	}

	return string(plain)
}

// Base64Content returns the base64 version of content.
func (l Wire) Base64Content() string {
	if l.Type == "base64" {
		return l.Content
	}

	return base64.StdEncoding.EncodeToString([]byte(l.Content))
}

// JSONStringContent returns the JSON content.
func (l Wire) JSONStringContent() string {
	if l.Type == "json" {
		return l.Content
	}

	if l.Type == "plain" {
		return "Error: Type is incorrect"
	}

	plain, err := base64.StdEncoding.DecodeString(l.Content)
	if err != nil {
		return "Failed to decode base64 content. Error: " + err.Error()
	}

	return string(plain)
}

// EncodePlain builds the wire protocol for plaintext type.
func (l Wire) EncodePlain() string {
	basic := fmt.Sprintf("type:plain|created:%v|content:%v", l.Created, l.PlainContent())

	if l.Topic != "" {
		return fmt.Sprintf("topic:%v|%v", l.Topic, basic)
	}

	return basic
}

// EncodePlain builds the wire protocol for base64 type.
func (l Wire) EncodeBase64() string {
	basic := fmt.Sprintf("type:base64|created:%v|content:%v", l.Created, l.Base64Content())

	if l.Topic != "" {
		return fmt.Sprintf("topic:%v|%v", l.Topic, basic)
	}

	return basic
}

func (l Wire) EncodeJSON() string {
	basic := fmt.Sprintf("type:json|created:%v|content:%v", l.Created, l.JSONStringContent())

	if l.Topic != "" {
		return fmt.Sprintf("topic:%v|%v", l.Topic, basic)
	}

	return basic
}
