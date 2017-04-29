package aredis

import (
	"services/assert"

	"context"

	"services/initializer"

	"github.com/Sirupsen/logrus"
	redis "gopkg.in/redis.v5"
)

var (
	// Client the actual pool to use with redis
	Client *redis.Client
)

type initRedis struct {
}

// Initialize try to create a redis pool
func (initRedis) Initialize(ctx context.Context) {
	Client = redis.NewClient(
		&redis.Options{
			Network:  cfg.Network,
			Addr:     cfg.Address,
			Password: cfg.Password,
			PoolSize: cfg.PoolSize,
			DB:       cfg.DB,
		},
	)
	// PING the server to make sure every thing is fine
	assert.Nil(Client.Ping().Err())
	logrus.Debug("redis is ready.")
	go func() {
		c := ctx.Done()
		assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
		<-c
		assert.Nil(Client.Close())
		logrus.Debug("redis finalized.")
	}()
}

func init() {
	initializer.Register(&initRedis{}, 0)
}
