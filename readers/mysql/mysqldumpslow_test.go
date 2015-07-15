package mysql

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var testMysqlSlowData = "/tmp/mysql-slow.log"

func generateTestDataFromMysqlDumpSlow() error {
	data := `Time                 Id Command    Argument
# Time: 150221 20:12:56
# User@Host: root[root] @ localhost []  Id:     1
# Query_time: 10.023847  Lock_time: 5.000116 Rows_sent: 1  Rows_examined: 10
use slow_testing;
SET timestamp=1424578376;
SELECT * FROM users WHERE name = 'Jesse';
# Time: 150221 20:13:30
# User@Host: root[root] @ localhost []  Id:     1
# Query_time: 18.000293  Lock_time: 11.000166 Rows_sent: 1  Rows_examined: 10
SET timestamp=1424578410;
SELECT * FROM users WHERE name = 'Jesse';
`
	return ioutil.WriteFile(testMysqlSlowData, []byte(data), 0644)
}

func TestMysqlDumpSlowRun(t *testing.T) {
	generateTestDataFromMysqlDumpSlow()

	m := &MysqlDumpSlow{}
	m.Data = make(map[string][]DumpSlow)
	m.Options = "-t 1 -s at"
	m.FilePath = testMysqlSlowData

	err := m.Run()
	if strings.Contains(err.Error(), "connection refused") {
		t.Infof("Local MySQL is not running. Stop testing.")
		return
	}

	if err != nil {
		t.Errorf("Fetching mysqldumpslow data should always be successful. Error: %v", err)
	}

	if len(m.Data["mysqldumpslow"]) == 0 {
		jsonData, _ := m.ToJson()
		t.Errorf("Processlist data should never be empty. Data: %v", string(jsonData))
	}

	os.Remove(testMysqlSlowData)
}
