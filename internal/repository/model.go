package repository

import (
	"time"

	"github.com/jackc/pgx/v5"
)

// Db представляет структуру базы данных
type Db struct {
	db *pgx.Conn
}

// NewDb создает новый экземпляр Db с подключением к базе данных
func NewDb(db *pgx.Conn) *Db {
	return &Db{
		db: db,
	}
}

// User представляет структуру пользователя
type User struct {
	Name  string
	Email string
}

// UserGET представляет структуру для получения пользователя из базы данных
type UserGET struct {
	ID        int
	User      User
	Role      bool
	CreatedAt time.Time
}
