package file

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"services/assert"
	"services/collection"
	"services/config"
	"services/initializer"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/expand"
)

var (
	rootPath *string
	allStore = make(initSys)
	lock     = &sync.Mutex{}
)

type initSys map[string]*store

func (m initSys) Initialize(c context.Context) {
	done := c.Done()
	assert.NotNil(done, "[BUG] init context has no done channel")
	func() {
		<-done
		for i := range m {
			if err := m[i].f.Close(); err != nil {
				logrus.Error(err)
			}
		}
	}()
}

type store struct {
	lock sync.Mutex
	name string
	f    *os.File
	enc  *json.Encoder
}

func (c store) Name() string {
	return c.name
}

type jsonSchema struct {
	Time time.Time   `json:"time"`
	Data interface{} `json:"data"`
}

func (c *store) Save(in interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.enc.Encode(jsonSchema{Time: time.Now(), Data: in})
}

func newFileCollection(n string) collection.Store {
	lock.Lock()
	defer lock.Unlock()

	if old, ok := allStore[n]; ok {
		return old
	}

	c := &store{
		lock: sync.Mutex{},
		name: n,
	}
	var err error
	c.f, err = os.OpenFile(filepath.Join(*rootPath, n), os.O_APPEND|os.O_CREATE, 0666)
	assert.Nil(err)

	c.enc = json.NewEncoder(c.f)
	allStore[n] = c
	return c
}

func init() {
	pwd, err := expand.Pwd()
	assert.Nil(err)
	rootPath = config.RegisterString("services.collection.file.root", pwd)

	// Register for the initializer, since we need to close files
	initializer.Register(allStore)
	// register as a driver
	collection.Register(newFileCollection)
}
