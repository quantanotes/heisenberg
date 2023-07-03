package common

import (
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/rs/zerolog/log"
)

const (
	logCacheSize = 1024
	maxPool      = 1024
	applyTimeout = 500 * time.Millisecond
)

func InitRaft(id string, path string, addr string, fsm raft.FSM) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(id)
	ip, sport, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(sport)
	if err != nil {
		return nil, err
	}
	addr = net.JoinHostPort(ip, fmt.Sprint(port*2))
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	store, err := raftboltdb.NewBoltStore(filepath.Join(path, "stable.raft"))
	if err != nil {
		return nil, err
	}
	cache, err := raft.NewLogCache(logCacheSize, store)
	if err != nil {
		return nil, err
	}
	snaps, err := raft.NewFileSnapshotStore(path, 3, log.Logger)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(addr, tcpAddr, maxPool, 5*time.Second, log.Logger)
	if err != nil {
		return nil, err
	}
	r, err := raft.NewRaft(config, fsm, cache, store, snaps, transport)
	if err != nil {
		return nil, err
	}
	r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raft.ServerID(id),
				Address: transport.LocalAddr(),
			},
		},
	})
	return r, nil
}

func ApplyRaft(r *raft.Raft, payload any) (any, error) {
	if r.State() != raft.Leader {
		return nil, fmt.Errorf("node is not leader of raft group")
	}
	cmd, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	future := r.Apply(cmd, applyTimeout)
	if err := future.Error(); err != nil {
		return nil, err
	}
	return future.Response(), nil
}
