package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func (s *Server) registerRoutes() {
	s.router = gin.Default()

	// Register health handler
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

func (s *Server) Run(port string, timeout int64) error {
	s.registerRoutes()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	go gracefulShutdown(server, timeout)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("error starting server: ", err)
		return err
	}

	log.Println("server stopped")
	return nil
}

func gracefulShutdown(server *http.Server, timeout int64) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}
}