package mysql

import (
	"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	Id           int            `db:"Id"`
	User         string         `db:"User"`
	Host         string         `db:"Host"`
	Db           sql.NullString `db:"db"`
	Command      string         `db:"Command"`
	Time         int            `db:"Time"`
	State        string         `db:"State"`
	Info         string         `db:"Info"`
	RowsSent     int            `db:"Rows_sent"`
	RowsExamined int            `db:"Rows_examined"`
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

		err := scanRows(rows, &plist)
		if err == nil {
			if m.Data[plist.Host] == nil {
				m.Data[plist.Host] = make([]Processlist, 0)
			}
			m.Data[plist.Host] = append(m.Data[plist.Host], plist)
		} else {
			return err
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *MysqlProcesslist) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}

func scanRows(rows *sqlx.Rows, plist *Processlist) error {
	var info sql.NullString
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Percona has the Rows_examined and Rows_sent columns
	if columns[len(columns)-1] == "Rows_examined" {
		err := rows.Scan(&plist.Id, &plist.User, &plist.Host, &plist.Db, &plist.Command, &plist.Time, &plist.State, &info, &plist.RowsSent, &plist.RowsExamined)
		if err != nil {
			return err
		}
	} else {
		err := rows.Scan(&plist.Id, &plist.User, &plist.Host, &plist.Db, &plist.Command, &plist.Time, &plist.State, &info)
		if err != nil {
			return err
		}
		plist.RowsSent = 0
		plist.RowsSent = 0
	}

	plist.Info = info.String

	return nil
}
