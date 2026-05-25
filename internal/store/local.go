package store

import (
	"fmt"
	"sync"

	"github.com/pr33tpatel/testFS/internal/tier"
)

// LocalStore manages chunks across storage tiers on an individual node
type LocalStore struct {
	mu    sync.RWMutex
	tiers []tier.StorageTier

	// chunkIndex maps chunkID -> the StorageTier it was written to
	chunkIndex map[string]tier.TierType
}

func NewLocalStore(tiers ...tier.StorageTier) (*LocalStore, error) {
	if len(tiers) == 0 {
		return nil, fmt.Errorf("[store.NewLocalStore]: at least one tier is required")
	}
	return &LocalStore{
		tiers:      tiers,
		chunkIndex: make(map[string]tier.TierType),
	}, nil
}

// Put writes a chunk \
// FIXME: default write to the first tier (memory) for now, IrisFS will select based on I/O signals
func (ls *LocalStore) Put(chunkID string, data []byte) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if _, exists := ls.chunkIndex[chunkID]; exists {
		return fmt.Errorf("[store.LocalStore.Put]: chunk %q already exists", chunkID)
	}

	// FIXME: default write to the first tier (memory) for now, IrisFS will select based on I/O signals
	t := ls.tiers[0]
	if err := t.Write(chunkID, data); err != nil {
		return fmt.Errorf("[store.LocalStore.Put]: %w", err)
	}
	ls.chunkIndex[chunkID] = t.Type() // map the chunkID to the StorageTier in the chunkIndex map
	return nil
}

// Get retrieves a chunk by using the chunkID to find the StorageTier it lives on
func (ls *LocalStore) Get(chunkID string) ([]byte, error) {
	ls.mu.RLock()
	tierType, exists := ls.chunkIndex[chunkID]
	defer ls.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("[store.LocalStore.Get]: chunk %q not found", chunkID)
	}

	t, err := ls.findTier(tierType)
	if err != nil {
		return nil, fmt.Errorf("[store.LocalStore.Get]: %w", err)
	}

	data, err := t.Read(chunkID)
	if err != nil {
		return nil, fmt.Errorf("[store.LocalStore.Get]: %w", err)
	}

	return data, nil
}

// Delete removes a chunk from the StorageTier it lives on, deletes the chunkEntry from the chunkIndex
func (ls *LocalStore) Delete(chunkID string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	tierType, exists := ls.chunkIndex[chunkID]
	if !exists {
		return fmt.Errorf("[store.LocalStore.Delete]: chunk %q not found", chunkID)
	}

	t, err := ls.findTier(tierType)
	if err != nil {
		return fmt.Errorf("[store.LocalStore.Delete]: %w", err)
	}

	if err := t.Delete(chunkID); err != nil {
		return fmt.Errorf("[store.LocalStore.Delete]: %w", err)
	}

	delete(ls.chunkIndex, chunkID) // delete the chunkEntry from the chunkIndex
	return nil
}

// TierStats returns Stats from every StorageTier on this node. \
// These are what get reported to the central server in heartbeats
func (ls *LocalStore) TierStats() []tier.Stats {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	stats := make([]tier.Stats, len(ls.tiers))
	for i, t := range ls.tiers {
		stats[i] = t.Stats()
	}
	return stats
}

// findTier retuns the tier matching the given "TierType"
// NOTE: must be called with a read lock active
func (ls *LocalStore) findTier(tt tier.TierType) (tier.StorageTier, error) {
	for _, t := range ls.tiers {
		if t.Type() == tt {
			return t, nil
		}
	}

	return nil, fmt.Errorf("[store.LocalStore.findTier]: no tier of type %q registered", tt)
}
