package inMemory

import (
	"errors"
	"ozonLinkShortener/internal/memory"
	"sync"
)

type inMemory struct {
	Longs  map[string]memory.MemoryEntry
	Shorts map[string]memory.MemoryEntry
}

func (c *inMemory) FindEntryByLongLink(longLink string) (entry memory.MemoryEntry, err error) {
	longsLock.Lock()
	defer longsLock.Unlock()
	value, ok := c.Longs[longLink]
	if !ok {
		return entry, errors.New("No such entry")
	}
	return value, nil
}

func (c *inMemory) FindEntryByShortId(shortId string) (entry memory.MemoryEntry, err error) {
	shortsLock.Lock()
	defer shortsLock.Unlock()
	value, ok := c.Shorts[shortId]
	if !ok {
		return entry, errors.New("No such entry")
	}
	return value, nil
}

func (c *inMemory) AddEntry(m memory.MemoryEntry) (err error) {
	shortsLock.Lock()
	longsLock.Lock()
	defer shortsLock.Unlock()
	defer longsLock.Unlock()
	c.Longs[m.LongLink] = m
	c.Shorts[m.ShortId] = m
	return nil
}

type inMemoryDb struct {
	inMemoryStorage *inMemory
}

var cacheSingleInstance inMemoryDb
var shortsLock = sync.Mutex{}
var longsLock = sync.Mutex{}

func GetMemoryInstance() inMemoryDb {
	if cacheSingleInstance.inMemoryStorage == nil {
		shortsLock.Lock()
		longsLock.Lock()
		defer shortsLock.Unlock()
		defer longsLock.Unlock()
		cacheSingleInstance = inMemoryDb{
			inMemoryStorage: &inMemory{
				Longs:  make(map[string]memory.MemoryEntry),
				Shorts: make(map[string]memory.MemoryEntry),
			},
		}
	}
	return cacheSingleInstance
}

func (c inMemoryDb) AddEntry(entry memory.MemoryEntry) error {
	return c.inMemoryStorage.AddEntry(entry)
}

func (c inMemoryDb) GetEntryByShortId(shortId string) (memory.MemoryEntry, error) {
	return c.inMemoryStorage.FindEntryByShortId(shortId)
}

func (c inMemoryDb) GetEntryByLongLink(longLink string) (memory.MemoryEntry, error) {
	return c.inMemoryStorage.FindEntryByLongLink(longLink)
}

func (c inMemoryDb) Clear() error {
	shortsLock.Lock()
	longsLock.Lock()
	defer shortsLock.Unlock()
	defer longsLock.Unlock()
	c.inMemoryStorage.Longs = make(map[string]memory.MemoryEntry)
	c.inMemoryStorage.Shorts = make(map[string]memory.MemoryEntry)
	return nil
}
