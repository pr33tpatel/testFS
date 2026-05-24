package registry

import "time"

type NodeInfo struct {
	NodeID     string
	Addr       string
	LastSeen   time.Time
	QueueDepth int
	LatencyMs  float64
}

type NodeRegistry interface {
	Register(info NodeInfo) error
	Update(nodeID string, info NodeInfo) error
	GetAll() []NodeInfo
	Get(nodeID string) (NodeInfo, error)
}
