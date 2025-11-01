package metrics

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// WatchdogGroup identifies a logical component that owns a set of metrics.
type WatchdogGroup string

const (
	GroupNode     WatchdogGroup = "node"
	GroupFarmer   WatchdogGroup = "farmer"
	GroupTimelord WatchdogGroup = "timelord"
)

type watchdog struct {
	group     WatchdogGroup
	name      string
	threshold time.Duration
	last      int64
	triggered uint32
	gauge     prometheus.Gauge
}

// WatchdogSnapshot represents the current status of a watchdog for JSON responses.
type WatchdogSnapshot struct {
	Group            WatchdogGroup `json:"group"`
	Metric           string        `json:"metric"`
	LastUpdate       time.Time     `json:"lastUpdate"`
	SecondsSince     float64       `json:"secondsSince"`
	ThresholdSeconds float64       `json:"thresholdSeconds"`
	Triggered        bool          `json:"triggered"`
}

var (
	watchdogStatuses = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "archivas_metrics_watchdog_triggered",
			Help: "Watchdog status for exported metrics (1=stale, 0=healthy)",
		},
		[]string{"group", "metric"},
	)

	watchdogsMu sync.Mutex
	watchdogs   = map[WatchdogGroup][]*watchdog{}
	started     = map[WatchdogGroup]bool{}
)

func registerWatchdog(group WatchdogGroup, metric string, threshold time.Duration) *watchdog {
	wd := &watchdog{
		group:     group,
		name:      metric,
		threshold: threshold,
		gauge:     watchdogStatuses.WithLabelValues(string(group), metric),
	}
	wd.Touch()

	watchdogsMu.Lock()
	watchdogs[group] = append(watchdogs[group], wd)
	watchdogsMu.Unlock()

	return wd
}

func (w *watchdog) Touch() {
	atomic.StoreInt64(&w.last, time.Now().UnixNano())
	atomic.StoreUint32(&w.triggered, 0)
	w.gauge.Set(0)
}

func (w *watchdog) start() {
	interval := w.threshold / 2
	if interval < 5*time.Second {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			lastNs := atomic.LoadInt64(&w.last)
			if lastNs == 0 {
				continue
			}

			last := time.Unix(0, lastNs)
			elapsed := time.Since(last)

			if elapsed > w.threshold {
				if atomic.SwapUint32(&w.triggered, 1) == 0 {
					log.Printf("[metrics] watchdog triggered: metric=%s stale=%s", w.name, elapsed)
				}
				w.gauge.Set(1)
			} else {
				if atomic.SwapUint32(&w.triggered, 0) == 1 {
					log.Printf("[metrics] watchdog recovered: metric=%s", w.name)
				}
				w.gauge.Set(0)
			}
		}
	}()
}

// StartWatchdogs launches watchdog monitoring loops for the provided group.
func StartWatchdogs(group WatchdogGroup) {
	watchdogsMu.Lock()
	if started[group] {
		watchdogsMu.Unlock()
		return
	}
	list := append([]*watchdog(nil), watchdogs[group]...)
	started[group] = true
	watchdogsMu.Unlock()

	for _, wd := range list {
		wd.start()
	}
}

// SnapshotWatchdogs returns a snapshot view of the current watchdog state for a group.
func SnapshotWatchdogs(group WatchdogGroup) []WatchdogSnapshot {
	watchdogsMu.Lock()
	list := append([]*watchdog(nil), watchdogs[group]...)
	watchdogsMu.Unlock()

	snapshots := make([]WatchdogSnapshot, 0, len(list))
	now := time.Now()
	for _, wd := range list {
		lastNs := atomic.LoadInt64(&wd.last)
		var last time.Time
		if lastNs > 0 {
			last = time.Unix(0, lastNs).UTC()
		}
		elapsed := now.Sub(last)
		if last.IsZero() {
			elapsed = 0
		}

		snapshots = append(snapshots, WatchdogSnapshot{
			Group:            wd.group,
			Metric:           wd.name,
			LastUpdate:       last,
			SecondsSince:     elapsed.Seconds(),
			ThresholdSeconds: wd.threshold.Seconds(),
			Triggered:        atomic.LoadUint32(&wd.triggered) == 1,
		})
	}

	return snapshots
}
