package alert

import (
	"site-monitor/internal/telegram"
	"site-monitor/pkg/logger"
	"time"

	"github.com/go-redis/redis"
	"github.com/segmentio/kafka-go"
)

type AlertConsumer struct {
	brokers  []string
	topic    string
	groupID  string
	reader   *kafka.Reader
	log      *logger.Logger
	telegram *telegram.Client
	redis    *redis.Client
	cooldown time.Duration
}

type AlertMessage struct {
	URL            string `json:"url"`
	Status         int    `json:"status"`
	ResponseTimeMs int    `json:"response_time_ms"`
	Error          string `json:"error"`
	Timestamp      string `json:"timestamp"`
}

type SiteState struct {
	IsUp      bool      `json:"is_up"`
	LastAlert time.Time `json:"last_alert"`
}
