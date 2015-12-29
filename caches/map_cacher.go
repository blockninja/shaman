package caches

// Simple cache that stores data in a simple go map.
// map doesn't automatically evict expired data, this will need to
// check to ensure data isn't already expired.

// TODO:
//  - add logging
//  - test
//  - add routine for removing old data

import (
	"time"

	"github.com/nanopack/shaman/config"
)

type mapCacher struct {
	expires int
	db      map[string]CacheEntry
}

// Map cacher initializer
func NewMapCacher(connection string, expires int) (*mapCacher, error) {
	config.Log.Info("creating map cacher")
	mc := mapCacher{expires: expires, db: make(map[string]CacheEntry)}
	return &mc, nil
}

func (self mapCacher) InitializeDatabase() error {
	return nil
}

func (self mapCacher) ClearDatabase() error {
	self.db = make(map[string]CacheEntry)
	return nil
}

// Get record from the map cacher and make sure it hasn't expired yet
func (self mapCacher) GetRecord(key string) (string, error) {
	var ce CacheEntry
	ce, ok := self.db[key]
	if !ok {
		config.Log.Debug("No Record: %s", key)
		return "", nil
	}
	if self.expires > 0 {
		if time.Now().Unix() > ce.Expires {
			// expired
			config.Log.Debug("Expired: %s", key)
			self.DeleteRecord(key)
			return "", nil
		}
		ce.Expires = time.Now().Unix() + int64(self.expires)
		self.db[key] = ce
	}
	config.Log.Debug("Found: %s = %s", key, ce.Value)
	return ce.Value, nil
}

// Insert/update entry in the map cacher
func (self mapCacher) SetRecord(key, val string) error {
	ce := CacheEntry{Expires: time.Now().Unix() + int64(self.expires), Value: val}
	self.db[key] = ce
	return nil
}

// Update entry in the map cacher
func (self mapCacher) ReviseRecord(key, val string) error {
	return self.SetRecord(key, val)
}

// remove entry from the map cacher
func (self mapCacher) DeleteRecord(key string) error {
	delete(self.db, key)
	return nil
}

func (self mapCacher) ListRecords() ([]string, error) {
	entries := make([]string, 0)
	for ce := range self.db {
		entries = append(entries, self.db[ce].Value)
	}
	return entries, nil
}
