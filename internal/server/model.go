package server

import (
	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/Ippolid/auth/postgres/query"
	"github.com/jackc/pgx/v4"
)

// Server представляет структуру сервера
type Server struct {
	auth_v1.UnimplementedAuthV1Server
	db *query.Db
}

// NewServer создает новый экземпляр сервера с подключением к базе данных
func NewServer(db *pgx.Conn) *Server {
	return &Server{
		db: query.NewDb(db),
	}
}
