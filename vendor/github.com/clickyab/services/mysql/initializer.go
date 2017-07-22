package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/healthz"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	gorp "gopkg.in/gorp.v2"
)

var (
	rdbmap   []*gorp.DbMap
	safeRead []*gorp.DbMap
	safeLock = sync.RWMutex{}
	wdbmap   *gorp.DbMap

	once    = sync.Once{}
	all     []Initializer
	factory func(string) (*sql.DB, error)
)

type initMysql struct {
}

type gorpLogger struct {
}

func (g gorpLogger) Printf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

func createDBMap(dsn, mark string) *gorp.DbMap {
	db, err := factory(dsn)
	assert.Nil(err)

	db.SetMaxIdleConns(maxIdleConnection.Int())
	db.SetMaxOpenConns(maxConnection.Int())

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	if develMode.Bool() {
		logger := gorpLogger{}
		dbmap.TraceOn(mark, logger)
	} else {
		dbmap.TraceOff()
	}

	return dbmap
}

func ping(db ...*gorp.DbMap) error {
	for i := range db {
		err := db[i].Db.Ping()
		if err != nil {
			logrus.Error(err)
			continue
		}
		return nil // just one active connection is fine
	}
	return fmt.Errorf("all %d ping(s) failed", len(db))
}

// Initialize the modules, its safe to call this as many time as you want.
func (in *initMysql) Initialize(ctx context.Context) {
	once.Do(func() {
		assert.NotNil(factory)

		wdbmap = createDBMap(wdsn.String(), "[wdb]")
		safe.Try(func() error { return ping(wdbmap) }, retryMax.Duration())

		rdsns := strings.Split(rdsnSlice.String(), ",")
		if len(rdsns) == 0 {
			logrus.Warn("No read db is configured. using write as read")
			rdsns = append(rdsns, wdsn.String()) // Add write as read source if there is no read source
		}

		for i := range rdsns {
			tmpDBMap := createDBMap(rdsns[i], fmt.Sprintf("[rdb%d]", i+1))
			rdbmap = append(rdbmap, tmpDBMap)
		}
		safe.Try(func() error { return ping(rdbmap...) }, retryMax.Duration())

		fillSafeArray()
		safe.GoRoutine(func() { updateRdbMap(ctx) })

		// Now that all initialization are done, lets initialize our modules
		for i := range all {
			all[i].Initialize()
		}

		healthz.Register(in)
		logrus.Debug("mysql is ready.")
		go func() {
			c := ctx.Done()
			assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
			<-c
			for _, i := range rdbmap {
				_ = i.Db.Close()
			}
			wdbmap.Db.Close()
			logrus.Debug("mysql finalized.")
		}()
	})
}

func updateRdbMap(ctx context.Context) {
	for {
		select {
		case <-time.After(rdbUpdateCD.Duration()):
			fillSafeArray()
		case <-ctx.Done():
			return
		}
	}
}

func fillSafeArray() {
	tmp := []*gorp.DbMap{}
	for i := range rdbmap {
		if err := ping(rdbmap[i]); err == nil {
			tmp = append(tmp, rdbmap[i])
		}
	}
	fail := len(tmp) == 0
	if fail { // heck! no read available? fallback to write
		tmp = append(tmp, wdbmap)
	}
	x := len(tmp)
	// simply return if there is no change. prevent useless lock
	if !fail && (len(rdbmap) == x && len(safeRead) == x) {
		return
	}
	safeLock.Lock()
	defer safeLock.Unlock()

	safeRead = tmp
}

// Healthy return true if the databases are ok and ready for ping
func (in *initMysql) Healthy(context.Context) error {
	rErr := ping(rdbmap...)
	wErr := ping(wdbmap)

	if rErr != nil || wErr != nil {
		return fmt.Errorf("mysql PING failed, read error was %s and write error was %s", rErr, wErr)
	}

	return nil
}

// Manager is a base manager for transaction model
type Manager struct {
	tx          *gorp.Transaction
	transaction bool
}

// InTransaction return true if this manager s in transaction
func (m *Manager) InTransaction() bool {
	return m.transaction
}

// Begin is for begin transaction
func (m *Manager) Begin() error {
	var err error
	if m.transaction {
		logrus.Panic("already in transaction")
	}
	m.tx, err = wdbmap.Begin()
	if err == nil {
		m.transaction = true
	}
	return err
}

// Commit is for committing transaction. panic if transaction is not started
func (m *Manager) Commit() error {
	if !m.transaction {
		logrus.Panic("not in transaction")
	}
	err := m.tx.Commit()
	if err != nil {
		return err
	}
	m.tx = nil
	m.transaction = false
	return nil
}

// Rollback is for RollBack transaction. panic if transaction is not started
func (m *Manager) Rollback() error {
	if !m.transaction {
		logrus.Panic("Not in transaction")
	}
	err := m.tx.Rollback()

	if err != nil {
		return err
	}

	m.transaction = false
	return nil
}

// GetRDbMap is for getting the current dbmap
func (m *Manager) GetRDbMap() gorp.SqlExecutor {
	if m.transaction {
		return m.tx
	}
	safeLock.RLock()
	defer safeLock.RUnlock()

	index := rand.Intn(len(safeRead))
	return safeRead[index]
}

// GetRSQLDB return the raw connection to database
func (m *Manager) GetRSQLDB() *sql.DB {
	safeLock.RLock()
	defer safeLock.RUnlock()

	index := rand.Intn(len(safeRead))
	return safeRead[index].Db
}

// GetWDbMap is for getting the current dbmap
func (m *Manager) GetWDbMap() gorp.SqlExecutor {
	if m.transaction {
		return m.tx
	}
	return wdbmap
}

// GetWSQLDB return the raw connection to database
func (m *Manager) GetWSQLDB() *sql.DB {
	return wdbmap.Db
}

// GetProperDBMap try to get the current writer for development mode
func (m *Manager) GetProperDBMap() gorp.SqlExecutor {
	if develMode.Bool() {
		return m.GetWDbMap()
	}
	return m.GetRDbMap()
}

// Hijack try to hijack into a transaction
func (m *Manager) Hijack(ts gorp.SqlExecutor) error {
	if m.transaction {
		return errors.New("already in transaction")
	}
	t, ok := ts.(*gorp.Transaction)
	if !ok {
		return errors.New("there is no transaction to hijack")
	}

	m.transaction = true
	m.tx = t

	return nil
}

// AddTable registers the given interface type with gorp. The table name
// will be given the name of the TypeOf(i).  You must call this function,
// or AddTableWithName, for any struct type you wish to persist with
// the given DbMap.
//
// This operation is idempotent. If i's type is already mapped, the
// existing *TableMap is returned
func (m *Manager) AddTable(i interface{}) *gorp.TableMap {
	return wdbmap.AddTable(i)
}

// AddTableWithName has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	return wdbmap.AddTableWithName(i, name)
}

// AddTableWithNameAndSchema has the same behavior as AddTable, but sets
// table.TableName to name.
func (m *Manager) AddTableWithNameAndSchema(i interface{}, schema string, name string) *gorp.TableMap {
	return wdbmap.AddTableWithNameAndSchema(i, schema, name)
}

// TruncateTables try to truncate tables , useful for tests
func (m *Manager) TruncateTables(tbl string) error {
	q := "TRUNCATE " + tbl
	_, err := wdbmap.Exec(q)
	return err
}

// Register a new initMysql module
func Register(m ...Initializer) {
	all = append(all, m...)
}

// RegisterConnectionFactory register a connection factory for sql connection.
func RegisterConnectionFactory(f func(string) (*sql.DB, error)) {
	factory = f
}

func init() {
	initializer.Register(&initMysql{}, 0)
}
