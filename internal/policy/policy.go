package policy

import (
	"time"

	"github.com/pr33tpatel/testFS/internal/registry"
)

// PlacementPolicy carries all signals that a policy might need to make a decision
// NOTE: policies may not use all fields, depends on policy implementation
type PlacementContext struct {
	Nodes       []registry.NodeInfo // contains QueueDepth and LatencyMs
	FileSize    int64
	Filename    string
	AccessCount int
	Timestamp   time.Time
}

// PlacementPolicy is the core interface every policy must implement
type PlacementPolicy interface {
	SelectNode(ctx PlacementContext) (registry.NodeInfo, error)
	Name() string
}

// PolicyOpts is passed to any policy constructor that need configuraiton,
// simple policies may ignore it, but complex/configurable policies use it for weights/params
type PolicyOpts struct {
	// Weights for composite policies (key = signal name, value = weight)
	Weights map[string]float64

	// Params for ML or threshold-based policies
	Params map[string]string
}
