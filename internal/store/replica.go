package store

type replica struct {
	shardId string
	clients *map[string]*StoreClient
}
