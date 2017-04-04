package main

import (
	"assert"
	"config"
	"fmt"
	"services/cache"
	_ "services/cache/redis"
	"services/initializer"
	_ "services/redis"
	"time"
)

type test struct {
	X int
	Y string
	Z float64
}

func main() {
	config.Initialize()
	defer initializer.Initialize()()

	t := test{
	//X: 11111,
	//Y: "s;wlswijdoiw",
	//Z: 1000.999,
	}

	cc := cache.CreateWrapper("SSS", &t)
	err := cache.Hit("SSS", cc)
	if err == nil {
		fmt.Println("cache hit")
		fmt.Printf("%+v", t)
	}

	err = cache.Do(cc, time.Hour, nil)
	assert.Nil(err)

}
