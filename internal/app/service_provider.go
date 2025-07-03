package app

import (
	"context"
	"log"

	"github.com/Ippolid/auth/internal/api/auth"

	"github.com/Ippolid/auth/internal/client/cache/redis"

	"github.com/Ippolid/auth/internal/api/user"
	"github.com/Ippolid/auth/internal/config"
	"github.com/Ippolid/auth/internal/repository"
	auth2 "github.com/Ippolid/auth/internal/repository/auth"
	redisCache "github.com/Ippolid/auth/internal/repository/redis"
	user2 "github.com/Ippolid/auth/internal/repository/user"
	"github.com/Ippolid/auth/internal/service"
	auth3 "github.com/Ippolid/auth/internal/service/auth"
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
	jwtConfig     config.JWTConfig
	accessConfig  config.AccessConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository
	authRepository repository.AuthRepository

	redisPool    *redigo.Pool
	serviceCache repository.CacheInterface
	redisClient  redis.Client

	userService service.UserService
	authService service.AuthService

	userController *user.Controller
	authController *auth.Controller

	//logger *zap.Logger
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

func (s *serviceProvider) GetAccessConfig(_ context.Context) config.AccessConfig {
	if s.accessConfig == nil {
		cfg, err := config.NewAccessConfig()
		if err != nil {
			log.Fatalf("failed to get Access config: %s", err.Error())
		}
		s.accessConfig = cfg
	}
	return s.accessConfig
}

func (s *serviceProvider) GetJWTConfig(_ context.Context) config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get TLS config: %s", err.Error())
		}
		s.jwtConfig = cfg
	}
	return s.jwtConfig
}

//func (s *serviceProvider) IntitLogger() *zap.Logger {
//	if s.logger == nil {
//		// Initialize the logger with a default configuration and InfoLevel.
//		l, err := logger.New("Info")
//		if err != nil {
//			log.Fatalf("failed to init logger: %v", err)
//		}
//		s.logger = l
//	}
//	return s.logger
//}

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

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = user2.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = auth2.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}
func (s *serviceProvider) GetCache(ctx context.Context) repository.CacheInterface {
	if s.serviceCache == nil {
		s.serviceCache = redisCache.NewRedisCache(s.GetRedisClient(ctx))
	}

	return s.serviceCache
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = user3.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
			s.GetCache(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = auth3.NewService(
			s.AuthRepository(ctx),
			s.TxManager(ctx),
			s.GetCache(ctx),
			s.GetJWTConfig(ctx),
			s.GetAccessConfig(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) UserController(ctx context.Context) *user.Controller {
	if s.userController == nil {
		s.userController = user.NewController(s.UserService(ctx))
	}

	return s.userController
}

func (s *serviceProvider) AuthController(ctx context.Context) *auth.Controller {
	if s.authController == nil {
		s.authController = auth.NewController(s.AuthService(ctx))
	}

	return s.authController
}
