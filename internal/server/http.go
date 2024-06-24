package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erizkiatama/gotu-assignment/internal/api/book"
	"github.com/erizkiatama/gotu-assignment/internal/api/order"
	"github.com/erizkiatama/gotu-assignment/internal/api/user"
	"github.com/erizkiatama/gotu-assignment/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	UserHandler  *user.Handler
	BookHandler  *book.Handler
	OrderHandler *order.Handler
}

func (s *Server) registerRoutes() {
	s.router = gin.Default()

	v1 := s.router.Group("/api/v1")

	// Register health handler
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Register user handler
	userGroup := v1.Group("/user")
	userGroup.POST("/register", s.UserHandler.Register)
	userGroup.POST("/login", s.UserHandler.Login)

	// Register book handler
	bookGroup := v1.Group("/book")
	bookGroup.GET("/", s.BookHandler.List)

	// Register order handler
	orderGroup := v1.Group("/order")
	orderGroup.POST("/", middleware.AuthorizeToken(), s.OrderHandler.CreateOrder)
	orderGroup.GET("/", middleware.AuthorizeToken(), s.OrderHandler.ListOrder)
	orderGroup.GET("/:order_id", middleware.AuthorizeToken(), s.OrderHandler.DetailOrder)
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
