package app

import (
	"fmt"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/db"
	"github.com/erizkiatama/gotu-assignment/internal/server"
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

	srv := server.Server{}

	return srv.Run(cfg.Server.Port, cfg.Server.ShutdownTimeMillis)
}
