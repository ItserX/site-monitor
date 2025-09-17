package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/segmentio/kafka-go"

	"site-monitor/internal/config"
	"site-monitor/pkg/logger"
	"site-monitor/pkg/metrics"
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

	cycleStart := time.Now()
	defer func() {
		metrics.CheckerCycleDuration.Observe(time.Since(cycleStart).Seconds())
		metrics.CheckerCycleTotal.Inc()
		c.pushMetricsToPrometheus()
	}()

	sites, err := c.fetchSitesFromAPI()
	if err != nil {
		metrics.CheckerAPIErrors.Inc()
		return
	}

	metrics.CheckerSitesProcessed.Set(float64(len(sites)))
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
	msg, err := json.Marshal(result)
	if err != nil {
		c.log.Sugar.Errorw("Failed to marshal JSON for Kafka", "url", result.URL, "error", err)
		return
	}

	err = c.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(result.URL),
			Value: []byte(msg),
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
	statusCode := "unknown"

	if err != nil {
		result.Success = false
		result.ErrorMsg = err.Error()
		statusCode = "error"

		metrics.SiteCheckErrors.WithLabelValues(url, "connection_error").Inc()
		metrics.SiteCheckSuccess.WithLabelValues(url).Set(0)

		c.log.Sugar.Errorw("Site check failed",
			"url", url,
			"error", err.Error(),
			"response_time_ms", result.ResponseTime,
		)

	} else {
		defer resp.Body.Close()
		result.StatusCode = resp.StatusCode
		statusCode = fmt.Sprintf("%d", resp.StatusCode)

		metrics.SiteCheckSuccess.WithLabelValues(url).Set(1)
		c.log.Sugar.Infow("Site check successed",
			"url", url,
			"status", result.StatusCode,
			"response_time_ms", result.ResponseTime,
		)
	}

	c.sendToKafka(result)

	metrics.SiteCheckTotal.WithLabelValues(url, statusCode).Inc()
	metrics.SiteCheckDuration.WithLabelValues(url, statusCode).Observe(float64(result.ResponseTime))

	return result
}

func (c *Checker) pushMetricsToPrometheus() {
	pusher := push.New(c.cfg.Prometheus.PushgatewayURL, "site_checker")

	pusher.Collector(metrics.SiteCheckTotal)
	pusher.Collector(metrics.SiteCheckSuccess)
	pusher.Collector(metrics.SiteCheckErrors)
	pusher.Collector(metrics.SiteCheckDuration)
	pusher.Collector(metrics.CheckerCycleTotal)
	pusher.Collector(metrics.CheckerCycleDuration)
	pusher.Collector(metrics.CheckerSitesProcessed)
	pusher.Collector(metrics.CheckerAPIErrors)

	if err := pusher.Push(); err != nil {
		c.log.Sugar.Errorw("Failed to push metrics to Pushgateway", "error", err)
	} else {
		c.log.Sugar.Debugw("Metrics pushed to Pushgateway successfully")
	}
}
