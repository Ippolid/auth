package query

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

func InsertUser(ctx context.Context, con *pgx.Conn, name, email, password string, role bool) (int, error) {
	var id int
	err := con.QueryRow(ctx,
		"INSERT INTO users_table (name,email,password,role) VALUES ($1,$2,$3,$4) RETURNING id", name, email, password, role).Scan(&id)

	if err != nil {

		return 0, err
	}

	return id, nil
}

type User struct {
	ID        int
	Name      string
	Email     string
	Role      bool
	CreatedAt time.Time
}

func GetUser(ctx context.Context, con *pgx.Conn, id int) (User, error) {
	row := con.QueryRow(ctx, "SELECT name,email,role,created_at FROM users_table WHERE id=$1", id)
	if row == nil {
		log.Fatalf("failed to get user: %v", row)
	}

	var user User
	var name, email string
	var role bool
	var createdAt time.Time
	if err := row.Scan(&name, &email, &role, &createdAt); err != nil {
		log.Fatalf("failed to scan user: %v", err)
		return User{}, err
	}
	user = User{
		ID:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: createdAt,
	}
	fmt.Println("111", user)

	return user, nil
}

func DeleteUser(ctx context.Context, con *pgx.Conn, id int) error {
	_, err := con.Exec(ctx, "DELETE FROM users_table WHERE id=$1", id)
	if err != nil {
		log.Fatalf("failed to delete note: %v", err)
	}
	return nil
}

func UpdateUser(ctx context.Context, con *pgx.Conn, id int, name, email string) error {
	_, err := con.Exec(ctx, "UPDATE users_table SET name=$1,email=$2 WHERE id=$3", name, email, id)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
	}

	return nil

}
