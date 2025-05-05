package app

import (
	"context"
	"log"

	"github.com/Ippolid/auth/internal/api/auth"
	"github.com/Ippolid/auth/internal/config"
	"github.com/Ippolid/auth/internal/repository"
	auth2 "github.com/Ippolid/auth/internal/repository/auth"
	"github.com/Ippolid/auth/internal/service"
	auth3 "github.com/Ippolid/auth/internal/service/auth"
	"github.com/Ippolid/platform_libary/pkg/closer"
	"github.com/Ippolid/platform_libary/pkg/db"
	"github.com/Ippolid/platform_libary/pkg/db/pg"
	"github.com/Ippolid/platform_libary/pkg/db/transaction"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	noteRepository repository.AuthRepository

	noteService service.AuthService

	noteController *auth.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.noteRepository == nil {
		s.noteRepository = auth2.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.noteService == nil {
		s.noteService = auth3.NewService(
			s.AuthRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.noteService
}

func (s *serviceProvider) NoteController(ctx context.Context) *auth.Controller {
	if s.noteController == nil {
		s.noteController = auth.NewController(s.AuthService(ctx))
	}

	return s.noteController
}
