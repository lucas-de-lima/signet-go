package main

import (
	"context"
	"crypto/ed25519"
	"sync"
	"time"
)

type cacheEntry struct {
	key       ed25519.PublicKey
	expiresAt time.Time
}

type CachingKeyResolver struct {
	provider *SlowKeyProvider
	ttl      time.Duration
	cache    sync.Map // map[string]cacheEntry
}

func NewCachingKeyResolver(provider *SlowKeyProvider, ttl time.Duration) *CachingKeyResolver {
	return &CachingKeyResolver{
		provider: provider,
		ttl:      ttl,
	}
}

// Resolve implementa KeyResolverFunc com cache TTL seguro para concorrência
func (c *CachingKeyResolver) Resolve(ctx context.Context, kid string) (ed25519.PublicKey, error) {
	if val, ok := c.cache.Load(kid); ok {
		entry := val.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.key, nil // cache hit
		}
		c.cache.Delete(kid) // expirado
	}
	// cache miss: busca no provider lento
	key, err := c.provider.FetchKey(kid)
	if err != nil {
		return nil, err
	}
	c.cache.Store(kid, cacheEntry{key: key, expiresAt: time.Now().Add(c.ttl)})
	return key, nil
}
