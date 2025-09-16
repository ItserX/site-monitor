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

	checkerCfg, redisCfg, err := loadConfigs()
	if err != nil {
		log.Sugar.Errorw("Failed to load configs", "error", err)
		return
	}

	redisClient, err := setupRedis(redisCfg, log)
	if err != nil {
		log.Sugar.Errorw("Failed to setup Redis", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupGracefulShutdown(cancel, log)
	runCheckerLoop(ctx, checkerCfg, redisClient, log)
}

func setupLogger() (*logger.Logger, error) {
	return logger.SetupLogger()
}

func loadConfigs() (config.CheckerConfig, config.RedisConfig, error) {
	var checkerCfg config.CheckerConfig
	var redisCfg config.RedisConfig

	if err := config.LoadConfig("configs/checker.yaml", &checkerCfg); err != nil {
		return checkerCfg, redisCfg, err
	}
	if err := config.LoadConfig("configs/redis.yaml", &redisCfg); err != nil {
		return checkerCfg, redisCfg, err
	}
	return checkerCfg, redisCfg, nil
}

func setupRedis(redisCfg config.RedisConfig, log *logger.Logger) (*storage.RedisStorage, error) {
	redisClient, err := storage.NewRedisStorage(redisCfg.Redis.Addr, redisCfg.Redis.Password, redisCfg.Redis.DB)
	if err != nil {
		log.Sugar.Errorw("Redis connection failed", "error", err)
		return nil, err
	}
	return redisClient, nil
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

func runCheckerLoop(ctx context.Context, cfg config.CheckerConfig, redisClient *storage.RedisStorage, log *logger.Logger) {
	ticker := time.NewTicker(time.Duration(cfg.Checker.Interval) * time.Second)
	defer ticker.Stop()

	log.Sugar.Infow("Checker service started", "interval_sec", cfg.Checker.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Sugar.Infow("Checker service stopped gracefully")
			return

		case <-ticker.C:
			checkSites(ctx, cfg, redisClient, log)
		}
	}
}

func checkSites(ctx context.Context, cfg config.CheckerConfig, redisClient *storage.RedisStorage, log *logger.Logger) {
	targets, err := redisClient.GetSites(ctx)
	if err != nil {
		log.Sugar.Errorw("Failed to get sites from Redis", "error", err)
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
