package kvstore

import "sync"

// store is the per-node in-memory shard storage.
// It is intentionally unexported; external code interacts only through Cluster.
type store struct {
	data map[string][]byte
	mu   sync.RWMutex
}

// newStore creates a new in-memory store.
func newStore() *store {
	return &store{
		data: make(map[string][]byte),
	}
}

// get returns (value, true) for a present key or (nil, false) for absent key.
// Time complexity: O(1) map lookup.
func (s *store) get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// set stores a value; overwrites if key exists.
// Time complexity: O(1).
func (s *store) set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// delete removes a key; no-op if absent.
// Returns true if key existed, false otherwise.
func (s *store) delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return ok
}

// keys returns a snapshot of all keys in the store.
func (s *store) keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]string, 0, len(s.data))
	for k := range s.data {
		result = append(result, k)
	}
	return result
}

// size returns the number of stored keys.
func (s *store) size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
