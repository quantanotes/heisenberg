package store

type StoreServer struct {
	id        string
	store     store
	shard     shard
	replica   replica
	isMaster  bool
	isReplica bool
}

func (s *StoreServer) Put() {

}
