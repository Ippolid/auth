package app

import (
	"context"
	"log"

	"github.com/Ippolid/auth/internal/client/cache/redis"

	"github.com/Ippolid/auth/internal/api/user"
	"github.com/Ippolid/auth/internal/config"
	"github.com/Ippolid/auth/internal/repository"
	redisCache "github.com/Ippolid/auth/internal/repository/redis"
	user2 "github.com/Ippolid/auth/internal/repository/user"
	"github.com/Ippolid/auth/internal/service"
	user3 "github.com/Ippolid/auth/internal/service/user"
	"github.com/Ippolid/platform_libary/pkg/closer"
	"github.com/Ippolid/platform_libary/pkg/db"
	"github.com/Ippolid/platform_libary/pkg/db/pg"
	"github.com/Ippolid/platform_libary/pkg/db/transaction"
	redigo "github.com/gomodule/redigo/redis"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	redisConfig   config.RedisConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	tlsConfig     config.TLSConfig

	dbClient       db.Client
	txManager      db.TxManager
	noteRepository repository.UserRepository

	redisPool    *redigo.Pool
	serviceCache repository.CacheInterface
	redisClient  redis.Client

	noteService service.UserService

	noteController *user.Controller
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

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// GetRedisConfig получаем конфиг для redis.
func (s *serviceProvider) GetRedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatal("failed to load redis config: %w", err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) GetTLSConfig() config.TLSConfig {
	if s.tlsConfig == nil {
		cfg, err := config.NewTLSConfig()
		if err != nil {
			log.Fatalf("failed to get TLS config: %s", err.Error())
		}
		s.tlsConfig = cfg
	}
	return s.tlsConfig
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

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.GetRedisConfig().MaxIdle(),
			IdleTimeout: s.GetRedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.GetRedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) GetRedisClient(_ context.Context) redis.Client {
	if s.redisConfig == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.GetRedisConfig())
	}

	return s.redisClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.UserRepository {
	if s.noteRepository == nil {
		s.noteRepository = user2.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}
func (s *serviceProvider) GetCache(ctx context.Context) repository.CacheInterface {
	if s.serviceCache == nil {
		s.serviceCache = redisCache.NewRedisCache(s.GetRedisClient(ctx))
	}

	return s.serviceCache
}

func (s *serviceProvider) AuthService(ctx context.Context) service.UserService {
	if s.noteService == nil {
		s.noteService = user3.NewService(
			s.AuthRepository(ctx),
			s.TxManager(ctx),
			s.GetCache(ctx),
		)
	}

	return s.noteService
}

func (s *serviceProvider) NoteController(ctx context.Context) *user.Controller {
	if s.noteController == nil {
		s.noteController = user.NewController(s.AuthService(ctx))
	}

	return s.noteController
}
