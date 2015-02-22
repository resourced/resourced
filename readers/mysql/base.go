package mysql

import (
	"github.com/jmoiron/sqlx"
)

var connections map[string]*sqlx.DB
