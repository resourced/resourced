package mysql

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("MysqlProcesslist", NewMysqlProcesslist)
}

func NewMysqlProcesslist() readers.IReader {
	m := &MysqlProcesslist{}
	m.Data = make(map[string][]Processlist)

	return m
}

// MysqlProcesslist is a reader that fetch SHOW FULL PROCESSLIST data.
type MysqlProcesslist struct {
	Data map[string][]Processlist
	Base
}

type Processlist struct {
	Id      int            `db:"Id"`
	User    string         `db:"User"`
	Host    string         `db:"Host"`
	Db      sql.NullString `db:"db"`
	Command string         `db:"Command"`
	Time    int            `db:"Time"`
	State   string         `db:"State"`
	Info    string         `db:"Info"`
}

func (m *MysqlProcesslist) Run() error {
	err := m.initConnection()
	if err != nil {
		return err
	}

	connection := connections[m.HostAndPort]

	rows, err := connection.Queryx("SHOW FULL PROCESSLIST")
	if err != nil {
		return err
	}

	for rows.Next() {
		var plist Processlist

		err := rows.StructScan(&plist)
		if err == nil {
			if m.Data[plist.Host] == nil {
				m.Data[plist.Host] = make([]Processlist, 0)
			}
			m.Data[plist.Host] = append(m.Data[plist.Host], plist)
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *MysqlProcesslist) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
