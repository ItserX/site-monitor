package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"site-monitor/internal/checker"
	"site-monitor/internal/config"
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

	checkerCfg, pgCfg, err := loadConfigs()
	if err != nil {
		log.Sugar.Errorw("Failed to load configs", "error", err)
		return
	}

	pgClient, err := setupPostgres(pgCfg, log)
	if err != nil {
		log.Sugar.Errorw("Failed to setup Postgres", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupGracefulShutdown(cancel, log)
	runCheckerLoop(ctx, checkerCfg, pgClient, log)
}

func setupLogger() (*logger.Logger, error) {
	return logger.SetupLogger()
}

func loadConfigs() (config.CheckerConfig, config.PostgresConfig, error) {
	var checkerCfg config.CheckerConfig
	var pgCfg config.PostgresConfig

	if err := config.LoadConfig("configs/checker.yaml", &checkerCfg); err != nil {
		return checkerCfg, pgCfg, err
	}
	if err := config.LoadConfig("configs/postgres.yaml", &pgCfg); err != nil {
		return checkerCfg, pgCfg, err
	}
	return checkerCfg, pgCfg, nil
}

func setupPostgres(pgCfg config.PostgresConfig, log *logger.Logger) (*storage.PostgresStorage, error) {
	client, err := storage.NewPostgresStorage(pgCfg.Postgres.DSN)
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

func runCheckerLoop(ctx context.Context, cfg config.CheckerConfig, pgClient *storage.PostgresStorage, log *logger.Logger) {
	ticker := time.NewTicker(time.Duration(cfg.Checker.Interval) * time.Second)
	defer ticker.Stop()

	log.Sugar.Infow("Checker service started", "interval_sec", cfg.Checker.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Sugar.Infow("Checker service stopped gracefully")
			return

		case <-ticker.C:
			checkSites(ctx, cfg, pgClient, log)
		}
	}
}

func checkSites(ctx context.Context, cfg config.CheckerConfig, pgClient *storage.PostgresStorage, log *logger.Logger) {
	targets, err := pgClient.GetSites(ctx)
	if err != nil {
		log.Sugar.Errorw("Failed to get sites from Postgres", "error", err)
		return
	}

	log.Sugar.Infow("Start checking sites", "count", len(targets))

	var wg sync.WaitGroup
	wg.Add(len(targets))

	for _, target := range targets {
		go func(url string) {
			defer wg.Done()
			result := checker.CheckSite(url, log, cfg.Checker.Timeout)
			log.Sugar.Infow("Site check completed",
				"url", result.URL,
				"success", result.Success,
				"status", result.StatusCode,
				"response_time_ms", result.ResponseTime,
				"error", result.ErrorMsg,
			)
		}(target.URL)
	}

	wg.Wait()
	log.Sugar.Infow("Check cycle completed", "sites_checked", len(targets))
}
