package permission

import (
	"sync"

	"github.com/clickyab/services/assert"
	"github.com/sirupsen/logrus"
)

// UserScope is the permission level for a role
// @Enum {
// }
type UserScope string

// Token is the resource to check
type Token string

const (
	// ScopeSelf means the user him self, no additional parameter
	ScopeSelf UserScope = "self"
	// ScopeGlobal means the entire perm, no param is required
	ScopeGlobal UserScope = "global"
)

// Interface is the perm interface
type Interface interface {
	// HasPermString is the has perm check
	Has(scope UserScope, perm Token, d int64) (UserScope, bool)
	// HasPermStringOn is the has perm on check
	HasOn(perm Token, ownerID int64, parentIDs []int64, d int64, scopes ...UserScope) (UserScope, bool)
}

// InterfaceComplete is the complete version of the interface to use
type InterfaceComplete interface {
	Interface
	// GetID return the id of holder
	GetID() int64
	// GetCurrentToken return the current permission that this object is built on
	GetCurrentToken() Token
	// GetCurrentScope return the current scope for this object (maximum)
	GetCurrentScope() UserScope
}

type complete struct {
	inner Interface

	id    int64
	perm  Token
	scope UserScope
}

var (
	registeredPerms = make(map[Token]string)
	lock            = &sync.RWMutex{}
)

const (
	// God is the god perm
	God Token = "god"
)

// HasPermString is the has perm check
func (pc complete) Has(scope UserScope, perm Token, d int64) (UserScope, bool) {
	return pc.inner.Has(scope, perm, d)
}

// HasPermStringOn is the has perm on check
func (pc complete) HasOn(perm Token, ownerID int64, parentIDs []int64, d int64, scopes ...UserScope) (UserScope, bool) {
	return pc.inner.HasOn(perm, ownerID, parentIDs, d, scopes...)
}

// GetID return the id of holder
func (pc complete) GetID() int64 {
	return pc.id
}

// GetCurrentToken return the current permission that this object is built on
func (pc complete) GetCurrentToken() Token {
	return pc.perm
}

// GetCurrentScope return the current scope for this object (maximum)
func (pc complete) GetCurrentScope() UserScope {
	return pc.scope
}

// Register register a permission
func Register(perm Token, description string) {
	lock.Lock()
	defer lock.Unlock()

	registeredPerms[perm] = description
}

// Registered check if the permission is registered in system or not
// and just log it
// TODO : panic if not
func Registered(perm Token) {
	lock.RLock()
	defer lock.RUnlock()

	if _, ok := registeredPerms[perm]; !ok {
		logrus.Errorf("The permission is not registered %s", perm)
	}

}

// GetAll return the permission list in system
func GetAll() map[Token]string {
	lock.RLock()
	defer lock.RUnlock()

	return registeredPerms
}

// NewInterfaceComplete return a new object base on the minimum object
func NewInterfaceComplete(inner Interface, id int64, perm Token, scope UserScope, domainID int64) InterfaceComplete {
	s, ok := inner.Has(scope, perm, domainID)
	if !ok {
		s, ok = inner.Has(ScopeGlobal, God, domainID)
	}
	assert.True(ok, "[BUG] probably there is some thing wrong with code generation")
	pc := &complete{
		inner: inner,
		id:    id,
		perm:  perm,
		scope: s,
	}

	return pc
}

// Column is a single column in data tables
type Column struct {
	Data           string            `json:"data"`
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	Searchable     bool              `json:"searchable"`
	Sortable       bool              `json:"sortable"`
	Visible        bool              `json:"visible"`
	Filter         bool              `json:"filter"`
	Title          string            `json:"title"`
	FilterValidMap map[string]string `json:"filter_valid_map"`
}

// Columns is the columns in data tables
type Columns []Column

func init() {
	Register(God, "the god, can do anything")
}
