package checker

import (
	"time"

	"github.com/segmentio/kafka-go"

	"site-monitor/internal/config"
	"site-monitor/pkg/logger"
)

const workerCount = 25

type Checker struct {
	cfg         config.CheckerConfig
	apiURL      string
	log         *logger.Logger
	kafkaWriter *kafka.Writer
}

type SiteCheckResult struct {
	URL          string    `json:"url"`
	StatusCode   int       `json:"status"`
	ResponseTime int64     `json:"response_time_ms"`
	Success      bool      `json:"success"`
	Timestamp    time.Time `json:"timestamp"`
	ErrorMsg     string    `json:"error"`
}

type Site struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}
