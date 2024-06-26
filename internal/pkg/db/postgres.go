package db

import (
	"fmt"
	"log"

	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDatabase(config config.PostgresConfig) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
	)
	return sqlx.MustConnect("postgres", dsn)
}

func RunDBMigrations(config config.PostgresConfig, url string) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	log.Println("start migrating schema...")
	migration, err := migrate.New(url, dsn)
	if err != nil {
		panic(err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	log.Println("migration success")
	return nil
}
