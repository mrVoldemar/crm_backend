package app

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/closer"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/config"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/config/services"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/cors"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/handlers"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/logger"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/middleware"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/ratelimit"
	"github.com/mrVoldemar/crm_backend/services/api-gw/internal/recovery"

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

var (
	configPath     string
	servicesConfig string
)

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.StringVar(&servicesConfig, "services-config", "config/services.json", "path to services config file")
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
	// Load services configuration
	serviceConfig, err := loadServicesConfig(servicesConfig)
	if err != nil {
		return err
	}

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(a.serviceProvider.GRPCConfig().AuthServiceURL())

	// Create rate limiter (100 requests per minute)
	rateLimiter := ratelimit.NewRateLimiter(100, time.Minute)

	// Create API handler
	apiHandler, err := handlers.NewAPIHandler(authMiddleware, rateLimiter, serviceConfig)
	if err != nil {
		return err
	}

	// Create mux and register routes
	mux := http.NewServeMux()
	apiHandler.RegisterRoutes(mux)

	// Wrap with middleware
	var handler http.Handler = mux
	handler = recovery.RecoveryMiddleware(handler)
	handler = logger.LoggingMiddleware(handler)
	handler = cors.CORSHeaders(handler)

	// Create HTTP server
	a.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: handler,
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

func loadServicesConfig(filePath string) (*services.Config, error) {
	// Open config file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse JSON
	config := &services.Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
