package kv

import "github.com/clickyab/services/assert"

// Scanner is the scanner interface
type Scanner interface {
	// Next is called each time to get the next set of keys
	Next(int) ([]string, bool)
	// Pattern return the pattern for current scanner
	Pattern() string
}

// ScannerFactory is the factory of the scanner
type ScannerFactory func(string) Scanner

var (
	scannerFactory ScannerFactory
)

// NewScanner return a new scanner
func NewScanner(pattern string) Scanner {
	regLock.RLock()
	defer regLock.RUnlock()
	assert.NotNil(scannerFactory, "[BUG] factory is not registered")

	return scannerFactory(pattern)
}
