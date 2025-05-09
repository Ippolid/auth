package model

import (
	"time"
)

// User структура пользователя
type User struct {
	ID        int64     `db:"id" redis:"id"`
	User      UserInfo  `db:""`
	Role      bool      `db:"role"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

// UserInfo информация о пользователя
type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
}

// Log структура для хранения логов
type Log struct {
	Method    string    `db:"method_name"`
	CreatedAt time.Time `db:"created_at"`
	Ctx       string    `db:"ctx"`
}
