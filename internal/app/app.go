package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Ippolid/auth/internal/metric"
	"github.com/Ippolid/auth/internal/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Ippolid/auth/internal/api/middleware"
	"github.com/Ippolid/auth/internal/config"
	"github.com/Ippolid/auth/internal/interceptor"
	"github.com/Ippolid/auth/internal/logger"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/Ippolid/auth/pkg/user_v1"
	"github.com/Ippolid/platform_libary/pkg/closer"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	_ "github.com/Ippolid/auth/statik" //nolint
)

// App представляет собой основное приложение
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

// NewApp создает новое приложение
func NewApp(ctx context.Context) (*App, error) {
	// Инициализируем логгер в первую очередь
	logger.InitLocalLogger("Info")
	logger.Info("Creating new App...")
	tracing.Init(logger.Logger(), "auth")

	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	logger.Info("App created successfully")

	err = metric.Init(ctx)
	if err != nil {
		logger.Error("Failed to init metric.", zap.Error(err))
	}
	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			logger.Fatal("failed to run GRPC server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			logger.Fatal("failed to run HTTP server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			logger.Error("failed to run Swagger server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()

		err := runPrometheus()
		if err != nil {
			logger.Fatal("failed to run GRPC server", zap.Error(err))
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load("./.env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	creds, err := credentials.NewServerTLSFromFile(
		a.serviceProvider.GetTLSConfig().CertFile(),
		a.serviceProvider.GetTLSConfig().KeyFile(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to load TLS credentials")
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.ServerTracingInterceptor,
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
				interceptor.MetricsInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)
	user_v1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserController(ctx))
	auth_v1.RegisterAuthServer(a.grpcServer, a.serviceProvider.AuthController(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	creds, err := credentials.NewClientTLSFromFile("/server_cert.pem", "")
	if err != nil {
		return errors.Wrap(err, "could not load client TLS credentials from service.pem")
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	grpcAddr := a.serviceProvider.GRPCConfig().Address()
	if err = user_v1.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcAddr, dialOpts); err != nil {
		return errors.Wrap(err, "failed to register UserV1 handler with grpc-gateway")
	}

	corsMiddleware := middleware.NewCorsMiddleware()
	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.SwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	address := a.serviceProvider.GRPCConfig().Address()
	logger.Info("GRPC server is running", zap.String("address", address))

	list, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	address := a.serviceProvider.HTTPConfig().Address()
	logger.Info("HTTP server is running", zap.String("address", address))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	address := a.serviceProvider.SwaggerConfig().Address()
	logger.Info("Swagger server is running", zap.String("address", address))

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

//nolint:revive
func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Serving swagger file", zap.String("path", path))

		statikFs, err := fs.New()
		if err != nil {
			logger.Error("Failed to init statik fs", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Opening swagger file", zap.String("path", path))
		file, err := statikFs.Open(path)
		if err != nil {
			logger.Error("Failed to open swagger file", zap.String("path", path), zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file http.File) {
			_ = file.Close()
		}(file)

		logger.Debug("Reading swagger file", zap.String("path", path))
		content, err := io.ReadAll(file)
		if err != nil {
			logger.Error("Failed to read swagger file", zap.String("path", path), zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			logger.Error("Failed to write swagger file to response", zap.String("path", path), zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Served swagger file successfully", zap.String("path", path))
	}
}

func runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:              ":2112",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Prometheus server is running on %s", "localhost:2112")

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
