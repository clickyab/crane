package aredis

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/healthz"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var (
	// Client the actual pool to use with redis
	Client RedisClient
	all    []initializer.Simple
	lock   sync.RWMutex
)

type initRedis struct {
}

// Healthy return true if the databases are ok and ready for ping
func (initRedis) Healthy(context.Context) error {
	ping, err := Client.Ping().Result()
	if err != nil || strings.ToUpper(ping) != "PONG" {
		return fmt.Errorf("redis PING failed. result was '%s' and the error was %s", ping, err)
	}

	return nil
}

// Initialize try to create a redis pool
func (i *initRedis) Initialize(ctx context.Context) {
	if cluster.Bool() {
		endpoints, err := lookup(address.String())
		assert.Nil(err)
		for i := range endpoints {
			endpoints[i] = fmt.Sprintf("%s:%d", endpoints[i], port.Int())
		}
		Client = redis.NewClusterClient(
			&redis.ClusterOptions{
				Addrs:    endpoints,
				Password: password.String(),
				PoolSize: poolsize.Int(),
			},
		)
	} else {
		Client = redis.NewClient(
			&redis.Options{
				Network:  "tcp",
				Addr:     fmt.Sprintf("%s:%d", address.String(), port.Int()),
				Password: password.String(),
				PoolSize: poolsize.Int(),
				DB:       db.Int(),
			},
		)
	}
	// PING the server to make sure every thing is fine
	safe.Try(func() error { return Client.Ping().Err() }, tryLimit.Duration())

	healthz.Register(i)

	for i := range all {
		all[i].Initialize()
	}
	logrus.Debug("redis is ready.")
	go func() {
		c := ctx.Done()
		assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
		<-c
		assert.Nil(Client.Close())
		logrus.Debug("redis finalized.")
	}()
}

// Register a new object to inform it after redis is loaded
func Register(in initializer.Simple) {
	lock.Lock()
	defer lock.Unlock()

	all = append(all, in)
}

func lookup(svcName string) ([]string, error) {
	var endpoints []string
	_, srvRecords, err := net.LookupSRV("", "", svcName)
	if err != nil {
		return endpoints, err
	}
	for _, srvRecord := range srvRecords {
		// The SRV records ends in a "." for the root domain
		ep := fmt.Sprintf("%v", srvRecord.Target[:len(srvRecord.Target)-1])
		endpoints = append(endpoints, ep)
	}
	fmt.Print(endpoints)
	return endpoints, nil
}

func init() {
	// Redis must be before mysql so the cache work on queries
	initializer.Register(&initRedis{}, -1)
}
