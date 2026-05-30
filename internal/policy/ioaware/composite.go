package ioaware

import (
	"fmt"
	"strconv"

	"github.com/pr33tpatel/testFS/internal/policy"
	"github.com/pr33tpatel/testFS/internal/registry"
)

func init() {
	policy.Register("ioaware-composite", func(opts policy.PolicyOpts) (policy.PlacementPolicy, error) {
		latW := weightOrDefault(opts.Weights, "latency", 0.5)
		queueW := weightOrDefault(opts.Weights, "queue", 0.5)
		return NewCompositePolicy(latW, queueW), nil
	})
}

// CompositePolicy scores nodes using a weighted combination of latency and queue depth.
// Weights are normalized so they don't need to sum to 1.0.
type CompositePolicy struct {
	latencyWeight float64
	queueWeight   float64
}

func NewCompositePolicy(latencyWeight, queueWeight float64) *CompositePolicy {
	return &CompositePolicy{
		latencyWeight: latencyWeight,
		queueWeight:   queueWeight,
	}
}

func (p *CompositePolicy) Name() string { return "ioaware-composite" }

func (p *CompositePolicy) SelectNode(ctx policy.PlacementContext) (registry.NodeInfo, error) {
	if len(ctx.Nodes) == 0 {
		return registry.NodeInfo{}, fmt.Errorf("policy.CompositePolicy.SelectNode: no nodes available")
	}

	bestScore := scoreNode(ctx.Nodes[0], p.latencyWeight, p.queueWeight)
	best := ctx.Nodes[0]

	for _, n := range ctx.Nodes[1:] {
		if s := scoreNode(n, p.latencyWeight, p.queueWeight); s < bestScore {
			bestScore = s
			best = n
		}
	}
	return best, nil
}

// scoreNode returns a lower-is-better composite score.
func scoreNode(n registry.NodeInfo, latW, queueW float64) float64 {
	return (latW * n.LatencyMs) + (queueW * float64(n.QueueDepth))
}

func weightOrDefault(weights map[string]float64, key string, def float64) float64 {
	if weights == nil {
		return def
	}
	if v, ok := weights[key]; ok {
		return v
	}
	_ = strconv.FormatFloat(def, 'f', 2, 64) // suppress unused import warning
	return def
}
