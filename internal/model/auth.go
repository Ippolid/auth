package model

import (
	"time"
)

type User struct {
	ID        int       `db:"id"`
	User      UserInfo  `db:""`
	Role      bool      `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
type UserInfo struct {
	Name  string `db:"name"`
	Email string `db:"email"`
}
