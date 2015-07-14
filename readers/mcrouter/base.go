// Package mcrouter gathers mcrouter related data from a host.
package mcrouter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type Base struct {
	HostAndPort string
	ConfigFile  string
}

func (mcr *Base) Stats() (map[string]interface{}, error) {
	stats, err := mcr.StatsFromNetcat()
	if err != nil {
		return nil, err
	}

	statsFromFile, err := mcr.StatsFromFile()
	if err != nil {
		return nil, err
	}

	for key, value := range statsFromFile {
		trimmedKey := strings.Replace(key, "libmcrouter.mcrouter.5000.", "", -1)
		stats[trimmedKey] = value
	}

	return stats, nil
}

func (mcr *Base) StatsFromNetcat() (map[string]interface{}, error) {
	addrParts := strings.Split(mcr.HostAndPort, ":")
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

func (mcr *Base) StatsFromFile() (map[string]interface{}, error) {
	addrParts := strings.Split(mcr.HostAndPort, ":")
	port := addrParts[1]

	statsFile := fmt.Sprintf("/var/mcrouter/stats/libmcrouter.mcrouter.%v.stats", port)

	statsJson, err := ioutil.ReadFile(statsFile)
	if err != nil {
		return nil, err
	}

	var stats map[string]interface{}
	err = json.Unmarshal(statsJson, &stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
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

		stats[statsKey] = value

		if statsKey == "pid" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "parent_pid" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "time" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "uptime" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "num_servers") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "num_clients" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "num_suspect_servers" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "mcc_txbuf_reqs" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "mcc_waiting_replies" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "destination_batch_size" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "asynclog_requests" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "proxy_reqs_processing" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "proxy_reqs_waiting" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "client_queue_notify_period" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "rusage") {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats[statsKey] = valueFloat64
			}
		}
		if statsKey == "ps_num_minor_faults" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "ps_num_major_faults" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "ps_user_time_sec" {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats[statsKey] = valueFloat64
			}
		}
		if statsKey == "ps_system_time_sec" {
			valueFloat64, err := strconv.ParseFloat(value, 64)
			if err == nil {
				stats[statsKey] = valueFloat64
			}
		}
		if statsKey == "ps_vsize" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "ps_rss" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if statsKey == "successful_client_connections" {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "fibers") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "cmd_cas_outlier") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "cmd_delete_outlier") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "cmd_get_outlier") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "cmd_gets_outlier") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "cmd_set_outlier") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
		if strings.HasPrefix(statsKey, "cmd_other_outlier") {
			valueInt64, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				stats[statsKey] = valueInt64
			}
		}
	}

	return stats
}
