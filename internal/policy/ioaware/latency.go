package ioaware

import (
	"fmt"

	"github.com/pr33tpatel/testFS/internal/policy"
	"github.com/pr33tpatel/testFS/internal/registry"
)

func init() {
	policy.Register("ioaware-latency", func(opts policy.PolicyOpts) (policy.PlacementPolicy, error) {
		return NewLatencyPolicy(), nil
	})
}

// LatencyPolicy selects the node with the lowest reported latency.
type LatencyPolicy struct{}

func NewLatencyPolicy() *LatencyPolicy { return &LatencyPolicy{} }

func (p *LatencyPolicy) Name() string { return "ioaware-latency" }

func (p *LatencyPolicy) SelectNode(ctx policy.PlacementContext) (registry.NodeInfo, error) {
	if len(ctx.Nodes) == 0 {
		return registry.NodeInfo{}, fmt.Errorf("[policy.LatencyPolicy.SelectNode]: no nodes available")
	}
	best := ctx.Nodes[0]
	for _, n := range ctx.Nodes[1:] {
		if n.LatencyMs < best.LatencyMs {
			best = n
		}
	}
	return best, nil
}
