package main

import (
	"context"
	"fmt"

	"site-monitor/internal/alert"
	"site-monitor/internal/config"
	"site-monitor/internal/telegram"
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

	alertCfg, err := loadConfig()
	if err != nil {
		log.Sugar.Errorw("Failed to load config", "error", err)
		return
	}

	tgClient := telegram.NewClient(alertCfg.Telegram.BotToken, alertCfg.Telegram.ChatID)
	consumer := alert.NewAlertConsumer(alertCfg, log, tgClient)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())

	utils.SetupGracefulShutdown(cancel, log)
	consumer.Consume(ctx)
}

func loadConfig() (config.AlertConfig, error) {
	var cfg config.AlertConfig
	if err := config.LoadConfig("configs/alert.yaml", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
