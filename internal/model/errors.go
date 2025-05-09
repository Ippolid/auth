package model

import (
	"errors"
)

var (
	// ErrUserNotFound нет пользователя в хранилище.
	ErrUserNotFound = errors.New("user not found")
)
