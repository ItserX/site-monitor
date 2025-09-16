package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"site-monitor/internal/config"
	"site-monitor/pkg/logger"
)

func NewChecker(cfg config.CheckerConfig, log *logger.Logger) *Checker {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Checker{
		cfg:         cfg,
		apiURL:      cfg.Checker.ApiURL,
		log:         log,
		kafkaWriter: writer,
	}
}

func (c *Checker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.cfg.Checker.Interval) * time.Second)
	defer ticker.Stop()

	c.log.Sugar.Infow("Checker service started", "interval_sec", c.cfg.Checker.Interval)

	for {
		select {
		case <-ctx.Done():
			c.log.Sugar.Infow("Checker service stopped gracefully")
			return
		case <-ticker.C:
			c.checkSites()
		}
	}
}

func (c *Checker) checkSites() {
	sites, err := c.fetchSitesFromAPI()
	if err != nil {
		return
	}

	c.log.Sugar.Infow("Start checking sites", "count", len(sites))

	jobs := make(chan Site, len(sites))
	results := make(chan SiteCheckResult, len(sites))

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for site := range jobs {
				results <- c.CheckSite(site.URL)
			}
		}()
	}

	for _, site := range sites {
		jobs <- site
	}
	close(jobs)

	wg.Wait()
	close(results)

	c.log.Sugar.Infow("Check cycle completed", "sites_checked", len(sites))
}

func (c *Checker) fetchSitesFromAPI() ([]Site, error) {
	client := http.Client{
		Timeout: time.Duration(c.cfg.Checker.Timeout) * time.Second,
	}

	resp, err := client.Get(c.apiURL + "/sites")
	if err != nil {
		c.log.Sugar.Errorw("Failed to fetch sites from API", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.log.Sugar.Errorw("API returned non-OK status", "status", resp.StatusCode)
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var sites []Site
	if err := json.NewDecoder(resp.Body).Decode(&sites); err != nil {
		c.log.Sugar.Errorw("Failed to decode API response", "error", err)
		return nil, err
	}

	return sites, nil
}

func (c *Checker) sendToKafka(result SiteCheckResult) {
	fmt.Println(123132)
	message := fmt.Sprintf(`{"url":"%s","status":%d,"response_time_ms":%d,"error":"%s","timestamp":"%s"}`,
		result.URL, result.StatusCode, result.ResponseTime, result.ErrorMsg, result.Timestamp.Format(time.RFC3339))

	err := c.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(result.URL),
			Value: []byte(message),
		},
	)
	if err != nil {
		c.log.Sugar.Errorw("Failed to send message to Kafka", "url", result.URL, "error", err)
	} else {
		c.log.Sugar.Infow("Send message to Kafka", "url", result)
	}
}

func (c *Checker) CheckSite(url string) SiteCheckResult {
	start := time.Now()
	result := SiteCheckResult{URL: url, Timestamp: start}

	client := http.Client{Timeout: time.Duration(c.cfg.Checker.Timeout) * time.Second}
	resp, err := client.Get(url)

	result.ResponseTime = time.Since(start).Milliseconds()

	if err != nil {
		result.Success = false
		result.ErrorMsg = err.Error()
		c.log.Sugar.Errorw("Site check failed",
			"url", url,
			"error", err.Error(),
			"response_time_ms", result.ResponseTime,
		)
		c.sendToKafka(result)

		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode < 400

	if result.Success {
		c.log.Sugar.Infow("Site check successed",
			"url", url,
			"status", result.StatusCode,
			"response_time_ms", result.ResponseTime,
		)
	}

	return result
}
