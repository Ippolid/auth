package query

import (
	"context"
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

func GetUsers(ctx context.Context, con *pgx.Conn, id int) ([]User, error) {
	rows, err := con.Query(ctx, "SELECT name,email,role,created_at FROM users_table WHERE id=$1", id)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		var id int
		var name, email string
		var role bool
		var createdAt time.Time

		err = rows.Scan(&id, &name, &email, &role, &createdAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}

		user = User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: createdAt,
		}

		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, name, email, createdAt)
		users = append(users, user)
	}

	return users, nil
}

func DeleteUser(ctx context.Context, con *pgx.Conn, id int) error {
	_, err := con.Exec(ctx, "DELETE FROM users_table WHERE id=$1", id)
	if err != nil {
		log.Fatalf("failed to delete note: %v", err)
	}

	return nil
}

func UpdateUser(ctx context.Context, con *pgx.Conn, name, email, password string, role bool) (int, error) {
	_, err := con.Exec(ctx, "UPDATE users_table SET name=$1,email=$2,password=$3,role=$4 WHERE id=$5", name, email, password, role)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
	}

	return 0, nil

}
