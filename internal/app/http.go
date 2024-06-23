package app

import (
	"github.com/erizkiatama/gotu-assignment/internal/config"
	"github.com/erizkiatama/gotu-assignment/internal/server"
)

func Initialize(cfg *config.Config) error {
	srv := server.Server{}

	return srv.Run(cfg.Server.Port, cfg.Server.ShutdownTimeMillis)
}
