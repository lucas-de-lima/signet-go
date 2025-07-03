package main

import (
	"crypto/ed25519"
	"fmt"
	"sync"
	"time"
)

type SlowKeyProvider struct {
	keys map[string]ed25519.PublicKey
	mu   sync.RWMutex
}

func NewSlowKeyProvider() *SlowKeyProvider {
	return &SlowKeyProvider{keys: make(map[string]ed25519.PublicKey)}
}

func (p *SlowKeyProvider) RegisterKey(kid string, pub ed25519.PublicKey) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.keys[kid] = pub
}

func (p *SlowKeyProvider) FetchKey(kid string) (ed25519.PublicKey, error) {
	time.Sleep(100 * time.Millisecond)
	p.mu.RLock()
	defer p.mu.RUnlock()
	key, ok := p.keys[kid]
	if !ok {
		return nil, fmt.Errorf("kid não encontrado: %s", kid)
	}
	return key, nil
}
