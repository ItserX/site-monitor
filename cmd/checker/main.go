package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"site-monitor/internal/checker"
	"site-monitor/internal/config"
	"site-monitor/pkg/logger"
)

func main() {
	log, err := setupLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		return
	}
	defer log.Sync()

	checkerCfg, err := loadConfigs()
	if err != nil {
		log.Sugar.Errorw("Failed to load config", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupGracefulShutdown(cancel, log)

	c := checker.NewChecker(checkerCfg, log)
	c.Run(ctx)
}

func setupLogger() (*logger.Logger, error) {
	return logger.SetupLogger()
}

func loadConfigs() (config.CheckerConfig, error) {
	var cfg config.CheckerConfig
	if err := config.LoadConfig("configs/checker.yaml", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
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
