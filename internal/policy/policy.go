package policy

import "github.com/pr33tpatel/testFS/internal/registry"

// PlacementPolicy determines which node should recieve a new chunk
type PlacementPolicy interface {
	SelectNode(nodes []registry.NodeInfo) (registry.NodeInfo, error)
}
