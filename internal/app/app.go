package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anton1ks96/college-app-core/internal/config"
	"github.com/anton1ks96/college-app-core/internal/handlers"
	"github.com/anton1ks96/college-app-core/internal/server"
	"github.com/anton1ks96/college-app-core/pkg/logger"
)

func Run() {
	cfg, err := config.Init()
	if err != nil {
		logger.Fatal(err)
	}

	handler := handlers.NewHandler(cfg)

	router := handler.Init()

	srv := server.NewServer(cfg, router)

	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Info(fmt.Sprintf("college-app-core started on port %s", cfg.Server.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logger.Error(fmt.Errorf("server forced to shutdown: %w", err))
	}

	logger.Info("server exited")
}
