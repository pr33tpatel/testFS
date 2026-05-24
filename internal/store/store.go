package store

import "github.com/pr33tpatel/testFS/internal/tier"

// ChunkStore manages chunks across different storage tiers on a node
type ChunkStore interface {
	Put(chunkID string, data []byte) error
	Get(chunkID string) ([]byte, error)
	Delete(chunkID string) error
	TierStats() []tier.Stats // reports all tiers back to server
}
