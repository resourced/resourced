package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var connections map[string]*sqlx.DB

func NewMysqlProcesslist() *MysqlProcesslist {
	m := &MysqlProcesslist{}
	m.Data = make(map[string][]Processlist)

	if connections == nil {
		connections = make(map[string]*sqlx.DB)
	}

	return m
}

// MysqlProcesslist is a reader that fetch SHOW FULL PROCESSLIST data.
type MysqlProcesslist struct {
	Data        map[string][]Processlist
	HostAndPort string
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

func (m *MysqlProcesslist) initConnection() error {
	var err error

	if _, ok := connections[m.HostAndPort]; !ok {
		connections[m.HostAndPort], err = sqlx.Open("mysql", fmt.Sprintf("root:@(%v)/?parseTime=true", m.HostAndPort))
	}

	return err
}

func (m *MysqlProcesslist) Run() error {
	err := m.initConnection()
	if err != nil {
		return err
	}

	connection := connections[m.HostAndPort]

	err = connection.Ping()
	if err != nil {
		return err
	}

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
