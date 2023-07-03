package common

type Meta = map[string]any

type Value struct {
	Index  uint64
	Vector []float32
	Meta   Meta
}

type ServiceType string

const (
	NoneService     ServiceType = "none"
	DatabaseService ServiceType = "database"
	StoreService    ServiceType = "store"
	CacheService    ServiceType = "cache"
)
