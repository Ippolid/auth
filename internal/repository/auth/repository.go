package auth

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/internal/utils"
	"github.com/Ippolid/platform_libary/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

const (
	tableName       = "users_table"
	tableAccessName = "access"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	createdAtColumn = "created_at"
	roleColumn      = "role"
	passwordColumn  = "password"
	tableLogName    = "logs"
	methodColumn    = "method_name"
	ctxColumn       = "ctx"
	endpointColumn  = "endpoint"
)

type repo struct {
	db db.Client
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

// InsertUser вставляет нового пользователя в базу данных и возвращает его ID

func (r *repo) Login(ctx context.Context, user model.LoginRequest) (*model.UserInfoJwt, error) {

	if user.Username == "" || user.Password == "" {
		return nil, fmt.Errorf("email and password must not be empty")
	}

	builder := sq.Select(passwordColumn, roleColumn).
		From(tableName).
		Where(sq.Eq{nameColumn: user.Username}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.Login",
		QueryRaw: query,
	}

	row := r.db.DB().QueryRowContext(ctx, q, args...)

	var userInfo model.UserInfoJwt
	var password string
	var role bool

	err = row.Scan(&password, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	if !utils.VerifyPassword(password, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}
	userInfo.Username = user.Username
	userInfo.Role = role

	return &userInfo, nil

}

func (r *repo) GetUserRole(ctx context.Context, username string) (bool, error) {
	builder := sq.Select(roleColumn).
		From(tableName).
		Where(sq.Eq{nameColumn: username}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.GetUserRole",
		QueryRaw: query,
	}

	row := r.db.DB().QueryRowContext(ctx, q, args...)

	var role bool
	err = row.Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("user not found")
		}
		return false, fmt.Errorf("failed to scan user role: %w", err)
	}

	return role, nil

}
func (r *repo) GetUsersAccess(ctx context.Context, isAdmin bool) ([]string, error) {
	builder := sq.Select(endpointColumn).
		From(tableAccessName).
		Where(sq.Eq{roleColumn: isAdmin}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.GetUsersAccess",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	endpoints := make([]string, 0)
	for rows.Next() {
		var endpoint string
		if err := rows.Scan(&endpoint); err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %w", err)
		}
		endpoints = append(endpoints, endpoint)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Если для указанной роли нет эндпоинтов, возвращаем пустой слайс
	if len(endpoints) == 0 {
		return endpoints, nil
	}

	return endpoints, nil
}

func (r *repo) MakeLog(ctx context.Context, info model.Log) error {
	builder := sq.Insert(tableLogName).
		Columns(methodColumn, createdAtColumn, ctxColumn).
		Values(info.Method, info.CreatedAt, info.Ctx).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.Log",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
