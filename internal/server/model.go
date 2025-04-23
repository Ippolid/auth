package server

import (
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/jackc/pgx/v5"
)

// Server представляет структуру сервера
type Server struct {
	auth_v1.UnimplementedAuthV1Server
	db *repository.Db
}

// NewServer создает новый экземпляр сервера с подключением к базе данных
func NewServer(db *pgx.Conn) *Server {
	return &Server{
		db: repository.NewDb(db),
	}
}
