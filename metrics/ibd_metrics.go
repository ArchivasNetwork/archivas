package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// IBD (Initial Block Download) metrics
// v1.1.1: Track sync progress and performance
var (
	// IBDRequestedBatches counts total batch requests sent
	IBDRequestedBatches = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_ibd_requested_batches_total",
		Help: "Total number of block batches requested during IBD",
	})

	// IBDReceivedBatches counts total batch responses received
	IBDReceivedBatches = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_ibd_received_batches_total",
		Help: "Total number of block batches received during IBD",
	})

	// IBDBlocksApplied counts total blocks applied during IBD
	IBDBlocksApplied = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_ibd_blocks_applied_total",
		Help: "Total number of blocks applied during IBD",
	})

	// IBDInflight tracks concurrent IBD streams
	IBDInflight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "archivas_ibd_inflight",
		Help: "Number of concurrent IBD streams",
	})

	// IBDBackoffSeconds tracks total time spent in backoff
	IBDBackoffSeconds = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_ibd_backoff_seconds_total",
		Help: "Total seconds spent in IBD backoff/retry",
	})
)

// Watchdogs for IBD metrics
var (
	ibdRequestedBatchesWatchdog = registerWatchdog(GroupNode, "archivas_ibd_requested_batches_total", 5*time.Minute)
	ibdReceivedBatchesWatchdog  = registerWatchdog(GroupNode, "archivas_ibd_received_batches_total", 5*time.Minute)
	ibdBlocksAppliedWatchdog    = registerWatchdog(GroupNode, "archivas_ibd_blocks_applied_total", 5*time.Minute)
)

// Helper functions to update IBD metrics and touch watchdogs
func IncIBDRequestedBatches() {
	IBDRequestedBatches.Inc()
	ibdRequestedBatchesWatchdog.Touch()
}

func IncIBDReceivedBatches() {
	IBDReceivedBatches.Inc()
	ibdReceivedBatchesWatchdog.Touch()
}

func IncIBDBlocksApplied(count int) {
	IBDBlocksApplied.Add(float64(count))
	ibdBlocksAppliedWatchdog.Touch()
}

func UpdateIBDInflight(value int) {
	IBDInflight.Set(float64(value))
}

func AddIBDBackoffSeconds(seconds float64) {
	IBDBackoffSeconds.Add(seconds)
}

