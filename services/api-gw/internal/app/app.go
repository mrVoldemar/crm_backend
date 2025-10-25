package app

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/closer"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/config"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/handlers"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	// Run both HTTP and gRPC servers
	go func() {
		if err := a.runHTTPServer(); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
		a.initGRPCServer,
	}
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(a.serviceProvider.GRPCConfig().AuthServiceURL())

	// Create API handler
	apiHandler, err := handlers.NewAPIHandler(authMiddleware)
	if err != nil {
		return err
	}

	// Create mux and register routes
	mux := http.NewServeMux()
	apiHandler.RegisterRoutes(mux)

	// Create HTTP server
	a.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer()

	reflection.Register(a.grpcServer)

	//userDesc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().GRPCAddress())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().GRPCAddress())
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
	log.Printf("HTTP server is running on %s", a.httpServer.Addr)

	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
