package ioaware

import (
	"fmt"

	"github.com/pr33tpatel/testFS/internal/policy"
	"github.com/pr33tpatel/testFS/internal/registry"
)

func init() {
	policy.Register("ioaware-queue", func(opts policy.PolicyOpts) (policy.PlacementPolicy, error) {
		return NewQueuePolicy(), nil
	})
}

// QueuePolicy selects the node with the lowest queue depth.
type QueuePolicy struct{}

func NewQueuePolicy() *QueuePolicy { return &QueuePolicy{} }

func (p *QueuePolicy) Name() string { return "ioaware-queue" }

func (p *QueuePolicy) SelectNode(ctx policy.PlacementContext) (registry.NodeInfo, error) {
	if len(ctx.Nodes) == 0 {
		return registry.NodeInfo{}, fmt.Errorf("policy.QueuePolicy.SelectNode: no nodes available")
	}
	best := ctx.Nodes[0]
	for _, n := range ctx.Nodes[1:] {
		if n.QueueDepth < best.QueueDepth {
			best = n
		}
	}
	return best, nil
}
