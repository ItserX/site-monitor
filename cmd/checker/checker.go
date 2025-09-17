package main

import (
	"context"
	"fmt"

	"site-monitor/internal/checker"
	"site-monitor/internal/config"
	"site-monitor/pkg/logger"
	"site-monitor/pkg/utils"
)

func main() {
	log, err := logger.SetupLogger()
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		return
	}
	defer log.Sync()

	checkerCfg, err := loadConfig()
	if err != nil {
		log.Sugar.Errorw("Failed to load config", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	utils.SetupGracefulShutdown(cancel, log)

	c := checker.NewChecker(checkerCfg, log)
	c.Run(ctx)
}

func loadConfig() (config.CheckerConfig, error) {
	var cfg config.CheckerConfig
	if err := config.LoadConfig("configs/checker.yaml", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
