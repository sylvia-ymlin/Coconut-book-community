package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP请求总数
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bookcommunity_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP请求延迟
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bookcommunity_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// 数据库查询延迟
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bookcommunity_db_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	// 缓存命中率
	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bookcommunity_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bookcommunity_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// 活跃用户数
	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bookcommunity_active_users",
			Help: "Number of currently active users",
		},
	)
)
