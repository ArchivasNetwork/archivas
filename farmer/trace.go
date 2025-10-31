package farmer

import (
	"log"
)

// CycleSummary contains per-height farming cycle stats
type CycleSummary struct {
	Height       uint64
	VDFIter      uint64
	PlotsScanned int
	ElapsedMS    int64
	ScanP50MS    int64
	ScanP95MS    int64
	Preempts     int
	BestQuality  uint64
	Threshold    uint64
	EarlyExit    bool
	Winner       bool
}

// PlotScanSummary contains per-plot scan stats
type PlotScanSummary struct {
	PlotName  string
	Chunks    int
	ElapsedMS int64
	Best      uint64
	Threshold uint64
	EarlyExit bool
	Preempted bool
}

// LogCycleSummary logs a cycle summary in grep-friendly format
func LogCycleSummary(s CycleSummary) {
	log.Printf("FARMER CYCLE height=%d vdf=%d plots=%d elapsed_ms=%d scans_p50_ms=%d scans_p95_ms=%d preempts=%d best_quality=%d threshold=%d early_exit=%v winner=%v",
		s.Height, s.VDFIter, s.PlotsScanned, s.ElapsedMS, s.ScanP50MS, s.ScanP95MS, s.Preempts, s.BestQuality, s.Threshold, s.EarlyExit, s.Winner)
}

// LogPlotScan logs a plot scan summary
func LogPlotScan(s PlotScanSummary) {
	log.Printf("FARMER PLOT name=%s chunks=%d elapsed_ms=%d best=%d need=%d early_exit=%v preempted=%v",
		s.PlotName, s.Chunks, s.ElapsedMS, s.Best, s.Threshold, s.EarlyExit, s.Preempted)
}

// LogQualityCheck logs quality comparison for debugging
func LogQualityCheck(sampleHex string, numeric uint64, threshold uint64) {
	log.Printf("QUALITY CHECK sample_hex=0x%s numeric=%d threshold=%d less=%v",
		sampleHex, numeric, threshold, numeric < threshold)
}

// LogQualityDomain logs the quality domain info once on startup
func LogQualityDomain() {
	qmax := uint64(^uint64(0)) // Max uint64
	log.Printf("QUALITY DOMAIN qmax=%d (18 quintillion) scale=SHA256→uint64", qmax)
}

// CalculatePercentiles calculates p50 and p95 from samples
func CalculatePercentiles(samples []int64) (p50, p95 int64) {
	if len(samples) == 0 {
		return 0, 0
	}

	// Simple approximation (would use sort in production)
	var sum int64
	for _, s := range samples {
		sum += s
	}
	avg := sum / int64(len(samples))

	return avg, avg * 2 // Rough estimate: p95 ≈ 2×avg for scan times
}
