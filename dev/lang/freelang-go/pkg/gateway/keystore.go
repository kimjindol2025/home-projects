package gateway

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// KeyStore manages API keys with CRUD operations.
type KeyStore struct {
	keys map[string]*APIKey
	mu   sync.RWMutex
}

// NewKeyStore creates a new empty key store.
func NewKeyStore() *KeyStore {
	return &KeyStore{
		keys: make(map[string]*APIKey),
	}
}

// Create generates a new API key with a random 32-byte secret.
func (ks *KeyStore) Create(keyID string, permissions []Permission, rateLimit int) (*APIKey, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	if _, exists := ks.keys[keyID]; exists {
		return nil, ErrDuplicateKey
	}

	// Generate 32-byte random secret
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	key := &APIKey{
		ID:          keyID,
		Secret:      secret,
		Permissions: permissions,
		RateLimit:   rateLimit,
		CreatedAt:   time.Now(),
		RevokedAt:   nil,
	}

	ks.keys[keyID] = key
	return key, nil
}

// Get retrieves an API key by ID.
func (ks *KeyStore) Get(keyID string) (*APIKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	key, exists := ks.keys[keyID]
	if !exists {
		return nil, ErrKeyNotFound
	}
	return key, nil
}

// Revoke marks a key as revoked by setting its RevokedAt timestamp.
func (ks *KeyStore) Revoke(keyID string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	key, exists := ks.keys[keyID]
	if !exists {
		return ErrKeyNotFound
	}

	if key.RevokedAt != nil {
		return fmt.Errorf("key already revoked")
	}

	now := time.Now()
	key.RevokedAt = &now
	return nil
}

// List returns all active (non-revoked) keys.
func (ks *KeyStore) List() []*APIKey {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	var active []*APIKey
	for _, key := range ks.keys {
		if key.IsActive() {
			active = append(active, key)
		}
	}
	return active
}

// Exists checks if a key ID exists and is active.
func (ks *KeyStore) Exists(keyID string) bool {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	key, exists := ks.keys[keyID]
	return exists && key.IsActive()
}
