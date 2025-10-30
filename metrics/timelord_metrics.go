package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Timelord metrics
var (
	// VDFIterationsTotal counts total VDF iterations computed
	VDFIterationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_vdf_iterations_total",
		Help: "Total VDF iterations computed",
	})

	// VDFTickDuration tracks time per VDF computation tick
	VDFTickDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "archivas_vdf_tick_duration_seconds",
		Help:    "Time per VDF computation tick",
		Buckets: prometheus.DefBuckets,
	})

	// VDFSeedChanged counts seed changes
	VDFSeedChanged = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_vdf_seed_changed_total",
		Help: "Total VDF seed changes (new blocks)",
	})
)

