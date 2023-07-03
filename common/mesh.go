package common

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/hashicorp/consul/api"
	"storj.io/drpc"
)

const configKey = "config"

type client struct {
	conn    drpc.Conn
	addr    string
	service ServiceType
	shard   uint
}

type meshConfig struct {
	NumCacheShards uint
}

type Mesh struct {
	consul  Consul
	clients map[string]*client
	config  meshConfig
	mu      sync.Mutex
}

func NewMesh(id string) (*Mesh, error) {
	consul, err := NewConsul(id)
	if err != nil {
		return nil, err
	}

	return &Mesh{
		consul:  *consul,
		clients: make(map[string]*client),
		config:  meshConfig{},
		mu:      sync.Mutex{},
	}, nil
}

func (m *Mesh) Init(ctx context.Context, addr string, stype ServiceType, shard uint) error {
	if err := m.consul.Register(addr, stype, shard); err != nil {
		return err
	}

	configData, err := m.consul.Get(configKey)
	if err != nil {
		return err
	}

	if configData == nil {
		m.config = meshConfig{NumCacheShards: 1}
		data, err := json.Marshal(&m.config)
		if err != nil {
			return err
		}
		m.consul.Put(configKey, data)
	} else if err := json.Unmarshal(configData, &m.config); err != nil {
		return err
	}

	if err := m.subscribe(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Mesh) Close() error {
	m.closeConns()
	return m.consul.Deregister()
}

func (m *Mesh) build(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.closeConns(); err != nil {
		return err
	}

	services, err := m.consul.Services()
	if err != nil {
		return err
	}

	m.clients = make(map[string]*client)
	for id, service := range services {
		conn, err := NewConn(ctx, service.Address)
		if err != nil {
			return nil
		}
		shard, err := strconv.ParseUint(service.Partition, 10, 64)
		if err != nil {
			return err
		}
		m.clients[id] = &client{
			conn:    conn,
			addr:    service.Address,
			service: ServiceType(service.Kind),
			shard:   uint(shard),
		}
	}

	return nil
}

func (m *Mesh) subscribe(ctx context.Context) error {
	if err := m.consul.Subscribe("services", "", func(val any) {
		m.build(ctx)
	}); err != nil {
		return err
	}

	if err := m.consul.Subscribe("key", "config", func(val any) {
		kv := val.(api.KVPair)
		if err := json.Unmarshal(kv.Value, &m.config); err != nil {
			return
		}
	}); err != nil {
		return err
	}

	return nil
}

func (m *Mesh) closeConns() error {
	for _, c := range m.clients {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
