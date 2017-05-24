package redis

import (
	"strconv"
	"time"

	"clickyab.com/exchange/services/redis"
	"clickyab.com/exchange/services/statistic"
)

type storeRedis struct {
	key    string
	expire time.Duration
}

func (sr *storeRedis) GetAll() (map[string]int64, error) {
	cmd := aredis.Client.HGetAll(sr.key)
	if err := cmd.Err(); err != nil {
		return nil, err
	}
	ss, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	si := make(map[string]int64, len(ss))
	for i := range ss {
		si[i], _ = strconv.ParseInt(ss[i], 10, 0)
	}

	return si, nil

}

func (sr *storeRedis) Key() string {
	return sr.key
}

func (sr *storeRedis) touchExpire() {
	// TODO : ignore the result for now. maybe later
	aredis.Client.Expire(sr.key, sr.expire)
}

func (sr *storeRedis) IncSubKey(s string, a int64) (int64, error) {
	cmd := aredis.Client.HIncrBy(sr.key, s, a)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	sr.touchExpire()
	return cmd.Val(), nil
}

func (sr *storeRedis) DecSubKey(s string, a int64) (int64, error) {
	cmd := aredis.Client.HIncrBy(sr.key, s, -a)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	sr.touchExpire()
	return cmd.Val(), nil
}

func (sr *storeRedis) Touch(s string) (int64, error) {
	cmd := aredis.Client.HGet(sr.key, s)
	if err := cmd.Err(); err != nil {
		return 0, err
	}
	sr.touchExpire()
	// this case is ok to ignore the err . just ignore it at touch
	i, _ := strconv.ParseInt(cmd.Val(), 10, 0)
	return i, nil
}

func factory(key string, expire time.Duration) statistic.Interface {
	return &storeRedis{
		key:    key,
		expire: expire,
	}
}

func init() {
	statistic.Register(factory)
}
