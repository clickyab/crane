package kv

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

// CacheProvider is the cacheFactory backend
type CacheProvider interface {
	// Do is called to store the cacheFactory
	Do(Cacheable, time.Duration) error
	// Hit called when we need to load the cacheFactory
	Hit(string, Cacheable) error
}

// CacheWrapper is a provider with support for inner entity
type CacheWrapper interface {
	Cacheable
	// Entity return the cached object
	Entity() interface{}
}

type cachable struct {
	entity interface{}
	key    string
}

var cacheFactory CacheProvider

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
	regLock.RLock()
	defer regLock.RUnlock()

	return cacheFactory.Do(e, t)
}

// Hit the cacheFactory
func Hit(key string, out Cacheable) error {
	regLock.RLock()
	defer regLock.RUnlock()

	return cacheFactory.Hit(key, out)
}

// CreateWrapper return an cachable object for this ntt
func CreateWrapper(key string, ntt interface{}) CacheWrapper {
	return &cachable{
		key:    key,
		entity: ntt,
	}
}
