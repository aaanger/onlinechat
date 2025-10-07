package main

import (
	"context"
	"log"
	"net/http"
	"onlineChat/internal/routes"
	"onlineChat/internal/users"
	"onlineChat/internal/ws"
	"onlineChat/pkg/config"
	"onlineChat/pkg/db"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
	logger     *logrus.Logger
}

func (srv *Server) Run(port string, handler http.Handler) error {
	srv.httpServer = &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	srv.logger.WithField("port", port).Info("Starting server")
	return srv.httpServer.ListenAndServe()
}

func (srv *Server) Shutdown() {
	srv.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.httpServer.Shutdown(ctx); err != nil {
		srv.logger.WithError(err).Error("Server forced to shutdown")
	} else {
		srv.logger.Info("Server exited gracefully")
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env file not found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger := setupLogger(cfg.Logging)
	logger.Info("Starting the server")

	gin.SetMode(cfg.Server.Mode)

	database, err := db.Open(cfg.Database)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping database")
	}
	logger.Info("Database connection established")

	userRepo := users.NewUserRepository(database, logger)
	chatRepo := ws.NewChatRepository(database, logger)

	userService := users.NewUserService(
		userRepo,
		cfg.JWT.Secret,
		time.Duration(cfg.JWT.ExpireHours)*time.Hour,
		logger,
	)

	chatService := ws.NewChatService(chatRepo, logger)

	hub := ws.NewHub(cfg.Redis, chatService, logger)
	go hub.Run()

	userHandler := users.NewUserHandler(userService, logger)
	chatHandler := ws.NewChatHandler(hub, chatService, logger)

	routeConfig := &routes.Config{
		JWT: routes.JWTConfig{
			Secret: cfg.JWT.Secret,
		},
		Security: routes.SecurityConfig{
			CORSOrigin:        cfg.Security.CORSOrigin,
			RateLimitRequests: cfg.Security.RateLimitRequests,
		},
	}

	router := routes.SetupRoutes(userHandler, chatHandler, routeConfig, logger)

	server := &Server{
		logger: logger,
	}

	go func() {
		if err := server.Run(cfg.Server.Port, router); err != nil {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.Shutdown()
}

func setupLogger(cfg config.LoggingConfig) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	return logger
}
