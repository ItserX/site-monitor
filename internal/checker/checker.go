package checker

import (
	"net/http"
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

func CheckSite(url string) SiteCheckResult {
	start := time.Now()
	result := SiteCheckResult{URL: url, Timestamp: time.Now()}

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)

	result.ResponseTime = time.Since(start).Milliseconds()

	if err != nil {
		result.Success = false
		result.ErrorMsg = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode < 400

	return result
}
