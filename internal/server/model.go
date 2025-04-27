package server

import (
	"context"
	"github.com/Ippolid/auth/internal/client/db/pg"
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/internal/repository/auth"
)

// Server представляет структуру сервера
type Server struct {
	serv *auth.Imp
	db   repository.AuthRepository
}

// NewServer создает новый экземпляр сервера с подключением к базе данных
func NewServer(dsn string) *Server {
	client, err := pg.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	return &Server{
		db: auth.NewRepository(client),
	}
}
