package policy

import (
	"fmt"
	"sync/atomic"

	"github.com/pr33tpatel/testFS/internal/registry"
)

func init() {
	Register("roundrobin", func(opts PolicyOpts) (PlacementPolicy, error) {
		return NewRoundRobin(), nil
	})
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

type RoundRobin struct {
	counter atomic.Uint64
}

func (rr *RoundRobin) Name() string {
	return "roundrobin"
}

func (rr *RoundRobin) SelectNode(ctx PlacementContext) (registry.NodeInfo, error) {
	if len(ctx.Nodes) == 0 {
		return registry.NodeInfo{}, fmt.Errorf("[policy.RoundRobin.SelectNode]: no nodes available")
	}
	idx := rr.counter.Add(1) - 1
	return ctx.Nodes[idx%uint64(len(ctx.Nodes))], nil
}
