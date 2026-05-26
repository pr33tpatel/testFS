package registry

import (
	"fmt"
	"sync"
	"time"
)

// NodeInfo defines properties of a node and metrics, NOTE: for testFS, metrics can be used for determining activity level
type NodeInfo struct {
	NodeID     string
	Addr       string
	LastSeen   time.Time
	QueueDepth int
	LatencyMs  float64
}

// NodeRegistry helps determine what nodes exists
type NodeRegistry interface {
	Register(info NodeInfo) error
	Update(nodeID string, info NodeInfo) error
	GetAll() []NodeInfo
	Get(nodeID string) (NodeInfo, error)
}

type MemoryNodeRegistry struct {
	mu    sync.RWMutex
	nodes map[string]NodeInfo
}

func NewMemoryNodeRegistry() *MemoryNodeRegistry {
	return &MemoryNodeRegistry{nodes: make(map[string]NodeInfo)}
}

func (r *MemoryNodeRegistry) Register(info NodeInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	info.LastSeen = time.Now()
	r.nodes[info.NodeID] = info
	return nil
}

func (r *MemoryNodeRegistry) Update(nodeID string, info NodeInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.nodes[nodeID]; !exists {
		return fmt.Errorf("registry.Update: node %q not registered", nodeID)
	}
	info.LastSeen = time.Now()
	r.nodes[nodeID] = info
	return nil
}

func (r *MemoryNodeRegistry) Get(nodeID string) (NodeInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.nodes[nodeID]
	if !ok {
		return NodeInfo{}, fmt.Errorf("registry.Get: node %q not found", nodeID)
	}
	return info, nil
}

func (r *MemoryNodeRegistry) GetAll() []NodeInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	nodes := make([]NodeInfo, 0, len(r.nodes))
	for _, n := range r.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}
