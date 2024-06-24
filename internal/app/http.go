package app

import (
	"fmt"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/erizkiatama/gotu-assignment/internal/pkg/db"
	"github.com/erizkiatama/gotu-assignment/internal/server"

	userApi "github.com/erizkiatama/gotu-assignment/internal/api/user"
	userRepository "github.com/erizkiatama/gotu-assignment/internal/repository/user"
	userService "github.com/erizkiatama/gotu-assignment/internal/service/user"

	bookApi "github.com/erizkiatama/gotu-assignment/internal/api/book"
	bookRepository "github.com/erizkiatama/gotu-assignment/internal/repository/book"
	bookService "github.com/erizkiatama/gotu-assignment/internal/service/book"

	orderApi "github.com/erizkiatama/gotu-assignment/internal/api/order"
	orderRepository "github.com/erizkiatama/gotu-assignment/internal/repository/order"
	orderService "github.com/erizkiatama/gotu-assignment/internal/service/order"
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
	bookRepo := bookRepository.New(database)
	orderRepo := orderRepository.New(database)

	// Initialize service
	userSvc := userService.New(userRepo)
	bookSvc := bookService.New(bookRepo)
	orderSvc := orderService.New(orderRepo)

	// Initialize handler
	userHandler := userApi.New(userSvc)
	bookHandler := bookApi.New(bookSvc)
	orderHandler := orderApi.New(orderSvc)

	srv := server.Server{
		UserHandler:  userHandler,
		BookHandler:  bookHandler,
		OrderHandler: orderHandler,
	}

	return srv.Run(cfg.Server.Port, cfg.Server.ShutdownTimeMillis)
}
