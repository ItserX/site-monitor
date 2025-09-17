package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	SiteCheckTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "site_check_total",
		Help: "Total number of site checks",
	}, []string{"url", "status"})

	SiteCheckDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "site_check_duration_ms",
		Help:    "Site check response time in milliseconds",
		Buckets: []float64{100, 200, 300, 500, 1000, 2000, 5000},
	}, []string{"url", "status"})

	SiteCheckSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "site_check_success",
		Help: "Site check success status (1 = success, 0 = failure)",
	}, []string{"url"})

	SiteCheckErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "site_check_errors_total",
		Help: "Total number of site check errors",
	}, []string{"url", "error_type"})

	CheckerCycleTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "checker_cycle_total",
		Help: "Total number of checker cycles",
	})

	CheckerSitesProcessed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "checker_sites_processed",
		Help: "Number of sites processed in last cycle",
	})

	CheckerCycleDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "checker_cycle_duration_seconds",
		Help:    "Duration of checker cycle in seconds",
		Buckets: []float64{1, 5, 10, 30, 60, 120},
	})

	CheckerAPIErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "checker_api_errors_total",
		Help: "Total number of API fetch errors",
	})
)
