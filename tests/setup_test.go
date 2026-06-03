package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alexisgalvan/onepiece-api/internal/config"
	"github.com/alexisgalvan/onepiece-api/internal/handler"
	postgresrepo "github.com/alexisgalvan/onepiece-api/internal/repository/postgres"
	redisrepo "github.com/alexisgalvan/onepiece-api/internal/repository/redis"
	"github.com/alexisgalvan/onepiece-api/internal/server"
	"github.com/alexisgalvan/onepiece-api/internal/usecase"
	"github.com/alexisgalvan/onepiece-api/pkg/database"
	"github.com/alexisgalvan/onepiece-api/pkg/logger"
	"go.uber.org/zap"
)

var (
	testApp *gin.Engine
	pgContainer *postgres.PostgresContainer
	redisContainer *redis.RedisContainer
)

func init() {
	logger.Init("test", "debug")
	gin.SetMode(gin.TestMode)
}

func setupTestEnv(ctx context.Context) (*gin.Engine, func(), error) {
	// Start Postgres
	pg, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
		postgres.WithInitScripts(filepath.Join("..", "migrations", "000001_initial_schema.up.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	pgHost, err := pg.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	pgPort, err := pg.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, nil, err
	}

	// Start Redis
	rd, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		return nil, nil, err
	}
	rdHost, err := rd.Host(ctx)
	if err != nil {
		return nil, nil, err
	}
	rdPort, err := rd.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, nil, err
	}

	// Configuration
	cfg := &config.Config{}
	cfg.App.Env = "test"
	
	cfg.Database.Host = pgHost
	cfg.Database.Port = pgPort.Port()
	cfg.Database.User = "test_user"
	cfg.Database.Password = "test_pass"
	cfg.Database.Name = "test_db"
	cfg.Database.SSLMode = "disable"

	cfg.Redis.Host = rdHost
	cfg.Redis.Port = rdPort.Port()

	// Database connections
	pgPool, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		return nil, nil, err
	}

	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, nil, err
	}

	// Repositories
	characterRepo := postgresrepo.NewCharacterRepository(pgPool)
	apiKeyRepo := postgresrepo.NewAPIKeyRepository(pgPool)
	cacheRepo := redisrepo.NewCacheRepository(redisClient)

	// Use cases
	characterUC := usecase.NewCharacterUseCase(characterRepo, cacheRepo)
	apiKeyUC := usecase.NewAPIKeyUseCase(apiKeyRepo)

	// Handlers
	healthHandler := handler.NewHealthHandler(cfg)
	characterHandler := handler.NewCharacterHandler(characterUC)
	apiKeyHandler := handler.NewAPIKeyHandler(apiKeyUC)

	// Setup Router
	router := server.SetupRouter(cfg, healthHandler, characterHandler, apiKeyHandler, apiKeyRepo, cacheRepo)

	cleanup := func() {
		pgPool.Close()
		redisClient.Close()
		pg.Terminate(context.Background())
		rd.Terminate(context.Background())
	}

	return router, cleanup, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	app, cleanup, err := setupTestEnv(ctx)
	if err != nil {
		logger.Log.Fatal("failed to setup test env", zap.Error(err))
	}
	testApp = app
	code := m.Run()
	cleanup()
	os.Exit(code)
}
