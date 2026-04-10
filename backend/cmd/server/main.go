package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/albal/disky/backend/internal/api/handlers"
	"github.com/albal/disky/backend/internal/config"
	"github.com/albal/disky/backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Load .env if exists
	_ = godotenv.Load("../../.env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer db.Close()

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "degraded", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	h := handlers.New(db)

	api := router.Group("/api")
	{
		api.GET("/products", h.GetProducts)
		
		auth := api.Group("/auth")
		{
			auth.GET("/google", h.GoogleLogin)
			auth.GET("/apple", h.AppleLogin)
			auth.GET("/microsoft", h.MicrosoftLogin)
		}

		alerts := api.Group("/alerts")
		{
			alerts.GET("/", h.GetAlerts)
			alerts.POST("/", h.CreateAlert)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Info().Msgf("server listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("listen error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server...")

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exiting")
}
