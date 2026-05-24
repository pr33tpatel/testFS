package tier

type TierType string

const (
	Memory TierType = "memory"
	Optane TierType = "optane"
	NVMe   TierType = "nvme"
	HDD    TierType = "hdd"
)

type Stats struct {
	QueueDepth int
	LatencyMs  float64
	FreeBytes  int64
	TierType   TierType
}

type StorageTier interface {
	Type() TierType
	Write(id string, data []byte) error
	Read(id string) ([]byte, error)
	Delete(id string) error
	Stats() Stats
}
