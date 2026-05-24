package tier

import (
	"fmt"
	"math/rand"
	"sync"
)

type MemoryTier struct {
	mu     sync.RWMutex
	chunks map[string][]byte
}

func NewMemoryTier() *MemoryTier {
	return &MemoryTier{
		chunks: make(map[string][]byte),
	}
}

func (m *MemoryTier) Type() TierType {
	return Memory
}

func (m *MemoryTier) Write(id string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	buf := make([]byte, len(data))
	copy(buf, data)
	m.chunks[id] = buf
	return nil
}

func (m *MemoryTier) Read(id string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.chunks[id]
	if !ok {
		return nil, fmt.Errorf("[tier.MemoryTier.Read]: chunk %q not found", id)
	}

	buf := make([]byte, len(data))
	copy(buf, data)
	return buf, nil
}

func (m *MemoryTier) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.chunks[id]; !ok {
		return fmt.Errorf("[tier.MemoryTier.Delete]: chunk %q not found", id)
	}

	delete(m.chunks, id)
	return nil
}

func (m *MemoryTier) Stats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	simQD := rand.Intn(20)                  // NOTE: Simulated Value, real hardware would read sysfs
	simLatencyMs := (rand.Float64()) * 1e-5 // NOTE: Simulated Value
	simFreeBytes := int64(1 << 30)          // NOTE: Simulated Value

	return Stats{
		TierType:   Memory,
		QueueDepth: simQD,
		LatencyMs:  simLatencyMs,
		FreeBytes:  simFreeBytes,
	}
}
