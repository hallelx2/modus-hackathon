package embedding

import (
	"sync"
	"time"
)

type EmbeddingCache struct {
	cache map[string]CacheEntry
	mu    sync.RWMutex
}

type CacheEntry struct {
	Embedding []float32
	Timestamp time.Time
	ExpiresAt time.Time
}

func NewEmbeddingCache() *EmbeddingCache {
	return &EmbeddingCache{
		cache: make(map[string]CacheEntry),
	}
}

func (ec *EmbeddingCache) Get(key string) ([]float32, bool) {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	entry, exists := ec.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Embedding, true
}

func (ec *EmbeddingCache) Set(key string, embedding []float32, expiresAt time.Time) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.cache[key] = CacheEntry{
		Embedding: embedding,
		Timestamp: time.Now(),
		ExpiresAt: expiresAt,
	}
}
