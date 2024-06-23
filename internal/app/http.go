package app

import "github.com/erizkiatama/gotu-assignment/internal/server"

func Initialize() error {
	srv := server.Server{}

	return srv.Run("8080", 5)
}
