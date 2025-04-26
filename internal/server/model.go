package server

import (
	"context"
	"github.com/Ippolid/auth/internal/client/db/pg"
	"github.com/Ippolid/auth/internal/repository"
	"github.com/Ippolid/auth/internal/repository/auth"
	"github.com/Ippolid/auth/pkg/auth_v1"
)

// Server представляет структуру сервера
type Server struct {
	auth_v1.UnimplementedAuthV1Server
	db repository.AuthRepository
}

// NewServer создает новый экземпляр сервера с подключением к базе данных
func NewServer(dsn string) *Server {
	client,err:=pg.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	return &Server{
		db: auth.NewRepository(client),
	}
}
