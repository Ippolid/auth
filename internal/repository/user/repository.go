package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ippolid/auth/internal/model"
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/platform_libary/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "users_table"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	createdAtColumn = "created_at"
	roleColumn      = "role"
	passwordColumn  = "password"
	tableLogName    = "logs"
	methodColumn    = "method_name"
	ctxColumn       = "ctx"
)

type repo struct {
	db db.Client
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

// InsertUser вставляет нового пользователя в базу данных и возвращает его ID
func (r *repo) CreateUser(ctx context.Context, user model.User) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(user.User.Name, user.User.Email, passwordHash, user.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "auth_repository.Create_User",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser получает пользователя по ID из базы данных
func (r *repo) GetUser(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(nameColumn, emailColumn, roleColumn, createdAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.Get_User",
		QueryRaw: query,
	}

	var user model.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}

	return &user, nil
}

// DeleteUser удаляет пользователя по ID из базы данных
func (r *repo) DeleteUser(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "auth_repository.Delete_User",
		QueryRaw: query,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	// Опционально: проверка, что запись действительно была удалена
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

// UpdateUser обновляет информацию о пользователе в базе данных
func (r *repo) UpdateUser(ctx context.Context, id int64, info model.UserInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, info.Name).
		Set(emailColumn, info.Email).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.Update_User",
		QueryRaw: query,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil

}

func (r *repo) MakeLog(ctx context.Context, info model.Log) error {
	builder := sq.Insert(tableLogName).
		PlaceholderFormat(sq.Dollar).
		Columns(methodColumn, createdAtColumn, ctxColumn).
		Values(info.Method, info.CreatedAt, info.Ctx)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Log",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
