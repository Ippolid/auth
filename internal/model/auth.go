package model

import (
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	User      UserInfo  `db:""`
	Role      bool      `db:"role"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
}
