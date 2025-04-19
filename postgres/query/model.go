package query

import (
	"github.com/jackc/pgx/v4"
	"time"
)

type Db struct {
	db *pgx.Conn
}

func NewDb(db *pgx.Conn) *Db {
	return &Db{
		db: db,
	}
}

type User struct {
	Name  string
	Email string
}

type UserGET struct {
	ID        int
	User      User
	Role      bool
	CreatedAt time.Time
}
