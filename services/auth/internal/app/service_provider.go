package app

import (
	"context"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/api/access"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/api/auth"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/service"
	"log"

	"github.com/mrVoldemar/crm_backend/services/auth/internal/client/db"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/client/db/pg"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/client/db/transaction"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/closer"
	"github.com/mrVoldemar/crm_backend/services/auth/internal/config"
	//"github.com/mrVoldemar/crm_backend/services/auth/internal/repository"
	//noteRepository "github.com/mrVoldemar/crm_backend/services/auth/internal/repository/note"
	//"github.com/mrVoldemar/crm_backend/services/auth/internal/service"
	accessService "github.com/mrVoldemar/crm_backend/services/auth/internal/service/access"
	authService "github.com/mrVoldemar/crm_backend/services/auth/internal/service/auth"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	jwtConfig  config.JwtConfig
	dbClient   db.Client
	txManager  db.TxManager
	//noteRepository repository.NoteRepository

	authService   service.AuthService
	accessService service.AccessService

	authImpl   *auth.Implementation
	accessImpl *access.Implementation
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
func (s *serviceProvider) JWTConfig() config.JwtConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJwtConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
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

/*
	func (s *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
		if s.noteRepository == nil {
			s.noteRepository = noteRepository.NewRepository(s.DBClient(ctx))
		}

		return s.noteRepository
	}
*/
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.JWTConfig(),
		)
	}

	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(
			s.JWTConfig(),
		)
	}

	return s.accessService
}
func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}
