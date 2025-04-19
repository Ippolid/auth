package server

import (
	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/Ippolid/auth/postgres/query"
	"github.com/jackc/pgx/v4"
)

type Server struct {
	auth_v1.UnimplementedAuthV1Server
	db *query.Db
}

func NewServer(db *pgx.Conn) *Server {
	return &Server{
		db: query.NewDb(db),
	}
}
