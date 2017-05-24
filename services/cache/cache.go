package cache

import (
	"encoding/gob"
	"io"
	"time"
)

// Serializable represent the object that can be serialized
type Serializable interface {
	// Decode is the decoder of this function
	Decode(io.Writer) error
	// Encode is the encoder function
	Encode(io.Reader) error
}

// Cacheable is the object that can be cached into
type Cacheable interface {
	Serializable

	String() string
}

// Provider is the cache backend
type Provider interface {
	// Do is called to store the cache
	Do(Cacheable, time.Duration) error
	// Hit called when we need to load the cache
	Hit(string, Cacheable) error
}

// Wrapper is a provider with support for inner entity
type Wrapper interface {
	Cacheable
	// Entity return the cached object
	Entity() interface{}
}

type cachable struct {
	entity interface{}
	key    string
}

var cache Provider

// Decode try to decode cookie profile into gob
func (cp *cachable) Decode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(cp.entity)
}

// Encode try to encode cookie profile from gob
func (cp *cachable) Encode(i io.Reader) error {
	dnc := gob.NewDecoder(i)
	return dnc.Decode(cp.entity)
}

func (cp *cachable) String() string {
	return cp.key
}

func (cp *cachable) Entity() interface{} {
	return cp.entity
}

// Do the entity
func Do(e Cacheable, t time.Duration, err error) error {
	if err != nil {
		return err
	}
	return cache.Do(e, t)
}

// Hit the cache
func Hit(key string, out Cacheable) error {
	return cache.Hit(key, out)
}

// CreateWrapper return an cachable object for this ntt
func CreateWrapper(key string, ntt interface{}) Wrapper {
	return &cachable{
		key:    key,
		entity: ntt,
	}
}

// Register a new cache provider
func Register(p Provider) {
	cache = p
}
