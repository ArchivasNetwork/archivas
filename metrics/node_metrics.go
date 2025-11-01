package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Node metrics - v1.0.1 restoration
var (
	// TipHeight tracks the current blockchain height
	TipHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "archivas_tip_height",
		Help: "Current blockchain tip height",
	})

	// PeerCount tracks number of connected peers
	PeerCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "archivas_peer_count",
		Help: "Number of connected P2P peers",
	})

	// BlocksSealed counts blocks sealed by this node
	BlocksSealed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_blocks_sealed_total",
		Help: "Total blocks sealed by this node",
	})

	// Difficulty tracks current mining difficulty
	Difficulty = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "archivas_difficulty",
		Help: "Current mining difficulty target in QMAX domain",
	})
	
	// Submission counters
	SubmitReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_submit_received_total",
		Help: "Proof submissions received",
	})
	
	SubmitAccepted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_submit_accepted_total",
		Help: "Proof submissions accepted",
	})
	
	SubmitIgnored = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_submit_ignored_total",
		Help: "Proof submissions ignored",
	})

	// BlockDuration tracks time to process blocks
	BlockDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "archivas_block_duration_seconds",
		Help:    "Time to process and apply a block",
		Buckets: prometheus.DefBuckets,
	})

	// IBDInflight tracks blocks being downloaded during IBD
	IBDInflight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "archivas_ibd_inflight",
		Help: "Number of blocks currently being downloaded",
	})

	// RPCRequests counts RPC requests by endpoint
	RPCRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "archivas_rpc_requests_total",
			Help: "Total RPC requests by endpoint",
		},
		[]string{"endpoint"},
	)

	// KnownPeers tracks total peers in peer store
	KnownPeers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "archivas_peers_known",
		Help: "Total number of known peers (connected + discovered)",
	})

	// GossipMessages counts peer gossip messages
	GossipMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_gossip_msgs_total",
		Help: "Total peer gossip messages sent",
	})

	// GossipAddrsReceived counts addresses received via gossip
	GossipAddrsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_gossip_addrs_received_total",
		Help: "Total peer addresses received via gossip",
	})

	// GossipDials counts auto-dial attempts from gossip
	GossipDials = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_gossip_dials_total",
		Help: "Total auto-dial attempts from gossip",
	})

	// GossipDialsFailed counts failed auto-dials
	GossipDialsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "archivas_gossip_dials_failed_total",
		Help: "Total failed auto-dial attempts",
	})
)

