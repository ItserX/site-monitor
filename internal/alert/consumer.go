package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/segmentio/kafka-go"

	"site-monitor/internal/config"
	"site-monitor/internal/telegram"
	"site-monitor/pkg/logger"
)

func NewAlertConsumer(cfg config.AlertConfig, log *logger.Logger, tg *telegram.Client) *AlertConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Kafka.Brokers,
		Topic:       cfg.Kafka.Topic,
		GroupID:     cfg.Kafka.GroupID,
		StartOffset: kafka.LastOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &AlertConsumer{
		brokers:  cfg.Kafka.Brokers,
		topic:    cfg.Kafka.Topic,
		groupID:  cfg.Kafka.GroupID,
		reader:   reader,
		log:      log,
		telegram: tg,
		redis:    rdb,
		cooldown: time.Duration(cfg.Redis.CooldownMinutes) * time.Minute,
	}
}

func (a *AlertConsumer) Consume(ctx context.Context) {
	a.log.Sugar.Infow("Listening for alerts on topic", "topic", a.topic)
	for {
		m, err := a.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				a.log.Sugar.Info("Shutting down consumer...")
				return
			}
			a.log.Sugar.Errorw("Error reading message", "error", err)
			continue
		}

		var alert AlertMessage
		if err := json.Unmarshal(m.Value, &alert); err != nil {
			a.log.Sugar.Errorw("Failed to parse alert JSON", "error", err, "raw", string(m.Value))
			continue
		}

		isUp := alert.Status == 200
		send, err := a.shouldSendAlert(alert.URL, isUp)
		if err != nil {
			a.log.Sugar.Errorw("Redis error", "error", err)
			continue
		}

		if !send {
			a.log.Sugar.Infow("No alert sent, status unchanged", "url", alert.URL, "status", alert.Status)
			continue
		}

		prettyMsg, _ := formatAlert(m.Value)
		if err := a.telegram.SendMessage(prettyMsg); err != nil {
			a.log.Sugar.Errorw("Failed to send Telegram alert", "error", err)
		} else {
			a.log.Sugar.Infow("Send alert to Telegram", "message", prettyMsg)
		}
	}
}

func (a *AlertConsumer) shouldSendAlert(url string, isUp bool) (bool, error) {
	key := "site_status:" + url
	val, err := a.redis.Get(key).Result()

	var state SiteState
	if err == nil {
		if err := json.Unmarshal([]byte(val), &state); err != nil {
			return true, nil
		}
	}

	now := time.Now()
	send := false

	if val == "" || state.IsUp != isUp {
		send = true
		state.IsUp = isUp
		state.LastAlert = now
	} else if now.Sub(state.LastAlert) > a.cooldown {
		send = true
		state.LastAlert = now
	}

	b, _ := json.Marshal(state)
	a.redis.Set(key, b, 0)

	return send, nil
}

func formatAlert(raw []byte) (string, error) {
	var alert AlertMessage
	if err := json.Unmarshal(raw, &alert); err != nil {
		return "", err
	}

	statusText := "❌ Unavailable"

	msg := fmt.Sprintf(
		"🚨 *Website Alert!*\n\n🌐 *URL*: %s\n📊 *Status*: %s (%d)\n⏱ *Response time*: %d ms\n🕒 *Timestamp*: %s",
		alert.URL,
		statusText,
		alert.Status,
		alert.ResponseTimeMs,
		alert.Timestamp,
	)

	if alert.Error != "" {
		msg += fmt.Sprintf("\n⚠️ *Error*: `%s`", alert.Error)
	}

	return msg, nil
}

func (a *AlertConsumer) Close() error {
	return a.reader.Close()
}
