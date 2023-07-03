package common

import (
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog/log"
)

type Consul struct {
	id     string
	client *api.Client
	plans  []*watch.Plan
}

func NewConsul(id string) (*Consul, error) {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Consul{id, client, nil}, nil
}

func (c *Consul) Register(hostport string, stype ServiceType, shard uint) error {
	_, sport, _ := net.SplitHostPort(hostport)
	port, _ := strconv.Atoi(sport)
	srv := api.AgentServiceRegistration{
		ID:        c.id,
		Name:      string(stype) + "-" + c.id,
		Port:      port,
		Kind:      api.ServiceKind(stype),
		Partition: strconv.FormatUint(uint64(shard), 10),
	}
	return c.client.Agent().ServiceRegister(&srv)
}

func (c *Consul) Deregister() error {
	return c.client.Agent().ServiceDeregister(c.id)
}

func (c *Consul) Services() (map[string]*api.AgentService, error) {
	return c.client.Agent().Services()
}

func (c *Consul) Subscribe(wtype string, key string, action func(val any)) error {
	params := map[string]any{
		"type": wtype,
		"key":  key,
	}
	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}

	plan.Handler = func(idx uint64, val interface{}) {
		action(val)
	}
	c.plans = append(c.plans, plan)

	go func() {
		err := plan.RunWithClientAndHclog(c.client, hclog.Default())
		if err != nil {
			log.Error().Err(err).Msg("Watch plan error")
		}
	}()

	return nil
}

func (c *Consul) Get(key string) ([]byte, error) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	return pair.Value, nil
}

func (c *Consul) Put(key string, value []byte) error {
	pair := &api.KVPair{Key: key, Value: value}
	_, err := c.client.KV().Put(pair, nil)
	return err
}

func (c *Consul) Delete(key string) error {
	_, err := c.client.KV().Delete(key, nil)
	return err
}
