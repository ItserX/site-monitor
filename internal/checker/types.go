package checker

import (
	"site-monitor/internal/config"
	"site-monitor/pkg/logger"
	"time"

	"github.com/segmentio/kafka-go"
)

const workerCount = 25

type Checker struct {
	cfg         config.CheckerConfig
	apiURL      string
	log         *logger.Logger
	kafkaWriter *kafka.Writer
}

type SiteCheckResult struct {
	URL          string
	StatusCode   int
	ResponseTime int64
	Success      bool
	Timestamp    time.Time
	ErrorMsg     string
}

type Site struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}
