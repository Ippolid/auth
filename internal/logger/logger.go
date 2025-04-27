package logger

import (
	"github.com/Ippolid/auth/internal/client/db"
)

const (
	tableName = "logs"

	idColumn        = "id"
	methodColumn    = "method_name"
	createdAtColumn = "created_at"
	errColumn       = "err"
)

type Logger struct {
	db db.Client
}

func NewLogger(db db.Client) Logger {
	return Logger{
		db: db,
	}
}
