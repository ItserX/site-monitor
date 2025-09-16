package checker

import (
	"net/http"
	"site-monitor/pkg/logger"
	"time"
)

type SiteCheckResult struct {
	URL          string
	StatusCode   int
	ResponseTime int64
	Success      bool
	Timestamp    time.Time
	ErrorMsg     string
}

func CheckSite(url string, log *logger.Logger, timeoutSeconds int) SiteCheckResult {
	start := time.Now()
	result := SiteCheckResult{URL: url, Timestamp: start}

	client := http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	resp, err := client.Get(url)

	result.ResponseTime = time.Since(start).Milliseconds()

	if err != nil {
		result.Success = false
		result.ErrorMsg = err.Error()
		log.Sugar.Errorw("Site check failed",
			"url", url,
			"error", err.Error(),
			"response_time_ms", result.ResponseTime,
		)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode < 400

	if result.Success {
		log.Sugar.Infow("Site check succeeded",
			"url", url,
			"status", result.StatusCode,
			"response_time_ms", result.ResponseTime,
		)
	} else {
		log.Sugar.Warnw("Site check returned error status",
			"url", url,
			"status", result.StatusCode,
			"response_time_ms", result.ResponseTime,
		)
	}

	return result
}
