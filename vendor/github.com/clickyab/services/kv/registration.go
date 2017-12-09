package kv

import "sync"

var (
	regLock sync.RWMutex
)

// Register is a gateway to register back end driver for all this types
func Register(eav KiwiFactory,
	store StoreFactory,
	dlock DistributedLockFactory,
	dset DistributedSetFactory,
	at StoreAtomicFactory,
	cache CacheProvider,
	scanner ScannerFactory,
	OneTime OneTimeFactory) {
	regLock.Lock()
	defer regLock.Unlock()

	kiwiFactory = eav
	storeFactory = store
	dlockFactory = dlock
	dsetFactory = dset
	atomicFactory = at
	cacheFactory = cache
	scannerFactory = scanner
	otFactory = OneTime
}
