package mysql

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/resourced/resourced/libprocess"
	"strconv"
	"strings"
)

func NewMysqlDumpSlow() *MysqlDumpSlow {
	m := &MysqlDumpSlow{}
	m.Data = make(map[string][]DumpSlow)

	return m
}

// MysqlDumpSlow is a reader that parses mysqldumpslow output.
type MysqlDumpSlow struct {
	Data     map[string][]DumpSlow
	Options  string
	FilePath string
}

type DumpSlow struct {
	Count     int
	User      string
	Time      string
	Lock      string
	Rows      float64
	Statement string
}

func (m *MysqlDumpSlow) validateBeforeRun() error {
	if m.Options == "" || m.FilePath == "" {
		return errors.New("Options or FilePath fields must not be empty.")
	}
	return nil
}

func (m *MysqlDumpSlow) Run() error {
	err := m.validateBeforeRun()
	if err != nil {
		return err
	}

	cmd := libprocess.NewCmd(fmt.Sprintf("mysqldumpslow %v %v", m.Options, m.FilePath))

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	m.Data["mysqldumpslow"] = make([]DumpSlow, 0)

	lineByLine := strings.Split(string(output), "\n")
	for i, line := range lineByLine {
		if strings.HasPrefix(line, "Count:") {
			// Example:
			// Count: 1  Time=0.00s (0s)  Lock=0.00s (0s)  Rows=1.0 (1), root[root]@localhost
			summaryAndUser := strings.Split(line, ",")
			summaryInChunk := strings.Fields(strings.TrimSpace(summaryAndUser[0]))

			d := DumpSlow{}
			d.Statement = strings.TrimSpace(lineByLine[i+1])
			d.User = strings.TrimSpace(summaryAndUser[1])

			count, err := strconv.Atoi(summaryInChunk[1])
			if err == nil {
				d.Count = count
			}

			timeAndValue := strings.Split(summaryInChunk[2], "=")
			d.Time = timeAndValue[1]

			lockAndValue := strings.Split(summaryInChunk[4], "=")
			d.Lock = lockAndValue[1]

			rowsAndValue := strings.Split(summaryInChunk[6], "=")
			rows, err := strconv.ParseFloat(rowsAndValue[1], 64)
			if err == nil {
				d.Rows = rows
			}

			m.Data["mysqldumpslow"] = append(m.Data["mysqldumpslow"], d)
		}
	}

	return err
}

// ToJson serialize Data field to JSON.
func (m *MysqlDumpSlow) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
