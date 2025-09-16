package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"site-monitor/internal/config"
	"site-monitor/internal/crud"
	"site-monitor/internal/storage"
	"site-monitor/pkg/logger"
)

func main() {
	log, err := setupLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		return
	}
	defer log.Sync()

	crudCfg, err := loadConfigs()
	if err != nil {
		log.Sugar.Errorw("Failed to load config", "error", err)
		return
	}

	pgClient, err := setupPostgres(crudCfg, log)
	if err != nil {
		log.Sugar.Errorw("Failed to setup Postgres", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupGracefulShutdown(cancel, log)

	runCrudServer(ctx, crudCfg, pgClient, log)
}

func setupLogger() (*logger.Logger, error) {
	return logger.SetupLogger()
}

func loadConfigs() (config.CrudConfig, error) {
	var crudCfg config.CrudConfig

	if err := config.LoadConfig("configs/crud.yaml", &crudCfg); err != nil {
		return crudCfg, err
	}

	return crudCfg, nil
}

func setupPostgres(crudCfg config.CrudConfig, log *logger.Logger) (storage.Storage, error) {
	var client storage.Storage
	client, err := storage.NewPostgresStorage(crudCfg.Postgres.DSN)
	if err != nil {
		log.Sugar.Errorw("Postgres connection failed", "error", err)
		return nil, err
	}
	return client, nil
}

func setupGracefulShutdown(cancelFunc context.CancelFunc, log *logger.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Sugar.Warnw("Shutdown signal received", "signal", sig.String())
		cancelFunc()
	}()
}

func runCrudServer(ctx context.Context, cfg config.CrudConfig, dbClient storage.Storage, log *logger.Logger) {
	r := chi.NewRouter()

	crudHandler := crud.NewHandler(dbClient, log)
	crudHandler.RegisterRoutes(r)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Sugar.Errorw("HTTP server Shutdown failed", "error", err)
		}
	}()

	log.Sugar.Infow("CRUD HTTP server started", "addr", srv.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Sugar.Errorw("HTTP server ListenAndServe failed", "error", err)
	}

	log.Sugar.Infow("CRUD HTTP server stopped gracefully")
}
