package store

type replica struct {
	clients *map[string]*StoreClient
}
