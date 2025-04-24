package repository

import "context"

type AuthRepository interface {
	InsertUser(ctx context.Context, user User, password string, role bool) (int, error)
	GetUser(ctx context.Context, id int) (UserGET, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, id int, user User) error
}
