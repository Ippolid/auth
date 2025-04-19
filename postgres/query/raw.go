package query

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

// InsertUser вставляет нового пользователя в базу данных и возвращает его ID
func (d *Db) InsertUser(ctx context.Context, user User, password string, role bool) (int, error) {
	var id int
	err := d.db.QueryRow(ctx,
		"INSERT INTO users_table (name,email,password,role) VALUES ($1,$2,$3,$4) RETURNING id", user.Name, user.Email, password, role).Scan(&id)

	if err != nil {

		return 0, err
	}

	return id, nil
}

// User представляет структуру пользователя

// GetUser получает пользователя по ID из базы данных
func (d *Db) GetUser(ctx context.Context, id int) (UserGET, error) {
	row := d.db.QueryRow(ctx, "SELECT name,email,role,created_at FROM users_table WHERE id=$1", id)

	var userg UserGET
	var user User
	var name, email string
	var role bool
	var createdAt time.Time
	if err := row.Scan(&name, &email, &role, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("no user found with id: %d", id)
			return UserGET{}, err
		}
		return UserGET{}, fmt.Errorf("failed to scan user: %w", err)
	}
	user = User{
		Name:  name,
		Email: email,
	}
	userg = UserGET{
		ID:        id,
		User:      user,
		Role:      role,
		CreatedAt: createdAt,
	}

	return userg, nil
}

// DeleteUser удаляет пользователя по ID из базы данных
func (d *Db) DeleteUser(ctx context.Context, id int) error {
	_, err := d.db.Exec(ctx, "DELETE FROM users_table WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// UpdateUser обновляет информацию о пользователе в базе данных
func (d *Db) UpdateUser(ctx context.Context, id int, user User) error {
	_, err := d.db.Exec(ctx, "UPDATE users_table SET name=$1,email=$2 WHERE id=$3", user.Name, user.Email, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil

}
