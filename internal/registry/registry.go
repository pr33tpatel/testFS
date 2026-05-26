package registry

import (
	"fmt"
	"sync"
)

type ChunkEntry struct {
	ChunkID  string
	NodeID   string
	Addr     string
	TierType string
}

// ChunkRegistry helps determine the location of each chunk
type ChunkRegistry interface {
	Assign(chunkID string, entry ChunkEntry) error
	Lookup(chunkID string) (ChunkEntry, error)
	Delete(chunkID string) error
}

type MemoryChunkRegistry struct {
	mu     sync.RWMutex
	chunks map[string]ChunkEntry
}

func NewMemoryChunkRegistry() *MemoryChunkRegistry {
	return &MemoryChunkRegistry{chunks: make(map[string]ChunkEntry)}
}

func (r *MemoryChunkRegistry) Assign(chunkID string, entry ChunkEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.chunks[chunkID]; exists {
		return fmt.Errorf("registry.Assign: chunk %q already assigned", chunkID)
	}
	r.chunks[chunkID] = entry
	return nil
}

func (r *MemoryChunkRegistry) Lookup(chunkID string) (ChunkEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.chunks[chunkID]
	if !ok {
		return ChunkEntry{}, fmt.Errorf("registry.Lookup: chunk %q not found", chunkID)
	}
	return entry, nil
}

func (r *MemoryChunkRegistry) Delete(chunkID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.chunks[chunkID]; !ok {
		return fmt.Errorf("registry.Delete: chunk %q not found", chunkID)
	}
	delete(r.chunks, chunkID)
	return nil
}
