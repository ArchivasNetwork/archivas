package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Farmer metrics
var (
	// PlotsLoaded tracks number of loaded plot files
	PlotsLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "archivas_plots_loaded",
		Help: "Number of plot files currently loaded",
	})

	// QualitiesChecked counts total quality checks
	QualitiesChecked = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_qualities_checked_total",
		Help: "Total number of plot quality checks performed",
	})

	// BlocksWonFarmer counts blocks successfully mined
	BlocksWonFarmer = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_blocks_won_total",
		Help: "Total number of blocks won by this farmer",
	})

	// FarmingLoopDuration tracks time per farming iteration
	FarmingLoopDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "archivas_farming_loop_seconds",
		Help:    "Time per farming loop iteration",
		Buckets: prometheus.DefBuckets,
	})
)

