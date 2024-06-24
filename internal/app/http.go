package app

import (
	"fmt"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/db"
	"github.com/erizkiatama/gotu-assignment/internal/server"

	userApi "github.com/erizkiatama/gotu-assignment/internal/api/user"
	userRepository "github.com/erizkiatama/gotu-assignment/internal/repository/user"
	userService "github.com/erizkiatama/gotu-assignment/internal/service/user"
)

func Initialize(cfg *config.Config) error {
	// Initialize database
	database := db.NewPostgresDatabase(cfg.Database.Postgres)
	defer func() {
		_ = database.Close()
	}()

	// Run database migrations
	if cfg.FeatureFlag.EnableMigrations {
		if err := db.RunDBMigrations(cfg.Database.Postgres, "file://scripts/migrations/"); err != nil {
			return fmt.Errorf("failed to migrate database: %v", err)
		}
	}

	// Initialize repository
	userRepo := userRepository.New(database)

	// Initialize service
	userSvc := userService.New(userRepo)

	// Initialize handler
	userHandler := userApi.New(userSvc)

	srv := server.Server{
		UserHandler: userHandler,
	}

	return srv.Run(cfg.Server.Port, cfg.Server.ShutdownTimeMillis)
}
