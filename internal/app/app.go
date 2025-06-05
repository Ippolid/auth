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

	"github.com/Ippolid/auth/pkg/auth_v1"

	"github.com/Ippolid/auth/internal/api/middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"

	"github.com/Ippolid/auth/internal/config"
	"github.com/Ippolid/auth/internal/interceptor"
	"github.com/Ippolid/auth/pkg/user_v1"
	"github.com/Ippolid/platform_libary/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
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
	log.Println("Creating new App...")
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}

	log.Println("App created successfully")
	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Printf("failed to run Swagger server: %v", err)
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
	// создаём TransportCredentials из файлов .crt и .key
	creds, err := credentials.NewServerTLSFromFile(
		a.serviceProvider.GetTLSConfig().CertFile(),
		a.serviceProvider.GetTLSConfig().KeyFile(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to load TLS credentials")
	}

	// создаём gRPC-сервер с TLS и middleware
	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	// включаем reflection и регистрируем сервис
	reflection.Register(a.grpcServer)
	user_v1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserController(ctx))
	auth_v1.RegisterAuthServer(a.grpcServer, a.serviceProvider.AuthController(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	// 2. Загружаем клиентские TLS-креденшлы из файла service.pem
	//    Второй аргумент "" означает, что Go возьмёт ServerName из адреса,
	//    который мы передадим в grpc.Dial (он должен совпадать с CN или SAN в сервисном сертификате).
	creds, err := credentials.NewClientTLSFromFile("/server_cert.pem", "")
	if err != nil {
		return errors.Wrap(err, "could not load client TLS credentials from service.pem")
	}

	// 3. Формируем grpc.DialOption с TLS-креденшлами
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// 4. Регистрируем grpc-gateway-handler, указываем адрес gRPC-сервера.
	//    address здесь — например, "localhost:50051" или как вернёт a.serviceProvider.GRPCConfig().Address().
	grpcAddr := a.serviceProvider.GRPCConfig().Address()
	if err := user_v1.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcAddr, dialOpts); err != nil {
		return errors.Wrap(err, "failed to register UserV1 handler with grpc-gateway")
	}

	// 5. Оборачиваем mux в CORS-мидлвар (если необходим) и создаём HTTP-сервер.
	corsMiddleware := middleware.NewCorsMiddleware()
	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(), // например, ":8080"
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
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
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
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

//nolint:revive
func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file http.File) {
			_ = file.Close()
		}(file)

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
