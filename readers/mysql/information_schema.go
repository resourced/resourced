package mysql

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlInformationSchemaTables() *MysqlInformationSchemaTables {
	m := &MysqlInformationSchemaTables{}
	m.Data = make(map[string][]InformationSchemaTables)

	return m
}

// MysqlInformationSchemaTables is a reader that fetch SHOW FULL PROCESSLIST data.
type MysqlInformationSchemaTables struct {
	Data map[string][]InformationSchemaTables
	Base
}

type InformationSchemaTables struct {
	Catalog              string        `db:"TABLE_CATALOG"`
	Schema               string        `db:"TABLE_SCHEMA"`
	Name                 string        `db:"TABLE_NAME"`
	Engine               string        `db:"ENGINE"`
	Version              int           `db:"VERSION"`
	RowFormat            string        `db:"ROW_FORMAT"`
	Rows                 uint64        `db:"TABLE_ROWS"`
	AvgRowLength         uint64        `db:"AVG_ROW_LENGTH"`
	DataLength           uint64        `db:"DATA_LENGTH"`
	MaxDataLength        uint64        `db:"MAX_DATA_LENGTH"`
	IndexLength          uint64        `db:"INDEX_LENGTH"`
	DataFree             uint64        `db:"DATA_FREE"`
	AutoIncrement        sql.NullInt64 `db:"AUTO_INCREMENT"`
	TableCollation       string        `db:"TABLE_COLLATION"`
	DataAndIndexLength   uint64
	DataAndIndexLengthMB float64
	DataAndIndexLengthGB float64
}

func (m *MysqlInformationSchemaTables) Run() error {
	err := m.initConnection()
	if err != nil {
		return err
	}

	query := `SELECT TABLE_CATALOG,TABLE_SCHEMA,TABLE_NAME,
ENGINE,VERSION,ROW_FORMAT,TABLE_ROWS,
AVG_ROW_LENGTH,DATA_LENGTH,MAX_DATA_LENGTH,INDEX_LENGTH,DATA_FREE,
AUTO_INCREMENT,TABLE_COLLATION FROM information_schema.TABLES ORDER BY (DATA_LENGTH + INDEX_LENGTH) DESC`

	rows, err := connections[m.HostAndPort].Queryx(query)
	if err != nil {
		return err
	}

	if m.Data["Tables"] == nil {
		m.Data["Tables"] = make([]InformationSchemaTables, 0)
	}

	for rows.Next() {
		var record InformationSchemaTables

		err := rows.StructScan(&record)
		if err == nil {
			record.DataAndIndexLength = record.DataLength + record.IndexLength
			record.DataAndIndexLengthMB = float64(record.DataAndIndexLength / 1024 / 1024)
			record.DataAndIndexLengthGB = float64(record.DataAndIndexLengthMB / 1024)

			m.Data["Tables"] = append(m.Data["Tables"], record)
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (m *MysqlInformationSchemaTables) ToJson() ([]byte, error) {
	return json.Marshal(m.Data)
}
