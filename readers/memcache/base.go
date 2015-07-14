// Package memcache gathers memcache related data from a host.
package memcache

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type Base struct {
	HostAndPort string
}

func (mc *Base) Stats() (map[string]interface{}, error) {
	addrParts := strings.Split(mc.HostAndPort, ":")
	host := addrParts[0]
	port := addrParts[1]

	if host == "" {
		host = "localhost"
	}

	c1 := exec.Command("echo", "stats")
	c2 := exec.Command("nc", host, port)

	var c2Output bytes.Buffer

	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = &c2Output

	c2.Start()

	err := c1.Run()
	if err != nil {
		return nil, err
	}

	err = c2.Wait()
	if err != nil {
		return nil, err
	}

	return NewStatsFromNetcat(c2Output.Bytes()), nil
}

func NewStatsFromNetcat(data []byte) map[string]interface{} {
	stats := make(map[string]interface{})

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		command := parts[0]

		if command != "STAT" {
			continue
		}

		statsKey := parts[1]
		value := strings.Join(parts[2:], " ")

		if statsKey == "version" || statsKey == "libevent" {
			stats[statsKey] = value

		} else if strings.HasPrefix(statsKey, "rusage") {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats[statsKey] = valueFloat64
			}

		} else {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
	}

	return stats
}
