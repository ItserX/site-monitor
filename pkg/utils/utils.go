package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"site-monitor/pkg/logger"
)

func SetupGracefulShutdown(cancelFunc context.CancelFunc, log *logger.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Sugar.Warnw("Shutdown signal received", "signal", sig.String())
		cancelFunc()
	}()
}
