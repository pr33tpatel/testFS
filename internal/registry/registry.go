package registry

// ChunkEntry tracks the location of chunks across nodes
type ChunkEntry struct {
	ChunkID  string
	NodeID   string
	Addr     string //host:port of the node
	TierType string
}

type ChunkRegistry interface {
	Assign(chunkID string, entry ChunkEntry) error
	Lookup(chunkID string) (ChunkEntry, error)
	Delete(chunkID string) error
}
