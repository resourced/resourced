package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMysqlProcesslist() *MysqlProcesslist {
	m := &MysqlProcesslist{}
	m.Data = make(map[string][]Processlist)
	return m
}

// MysqlProcesslist is a reader that fetch SHOW FULL PROCESSLIST data.
type MysqlProcesslist struct {
	Data        map[string][]Processlist
	HostAndPort string
}

type Processlist struct {
	Id      int            `db:"Id"`
	User    sql.NullString `db:"User"`
	Host    sql.NullString `db:"Host"`
	Db      sql.NullString `db:"db"`
	Command sql.NullString `db:"Command"`
	Time    int            `db:"Time"`
	State   sql.NullString `db:"State"`
	Info    sql.NullString `db:"Info"`
}

func (m *MysqlProcesslist) Run() error {
	sqlSession, err := sqlx.Open("mysql", fmt.Sprintf("root:@(%v)/?parseTime=true", m.HostAndPort))
	if err != nil {
		return err
	}

	err = sqlSession.Ping()
	if err != nil {
		return err
	}

	query := "SHOW FULL PROCESSLIST"
	rows, err := sqlSession.Queryx(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var plist Processlist

		err := rows.StructScan(&plist)
		if err == nil {
			if m.Data[plist.Host.String] == nil {
				m.Data[plist.Host.String] = make([]Processlist, 0)
			}
			m.Data[plist.Host.String] = append(m.Data[plist.Host.String], plist)
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *MysqlProcesslist) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
