package kvstore

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
)

const DefaultVirtualNodes = 150

// ringEntry represents one virtual node slot on the consistent hash ring.
type ringEntry struct {
	hash   uint32
	nodeID string
}

// Ring implements consistent hashing over a sorted slice of ringEntries.
type Ring struct {
	entries      []ringEntry
	nodeSet      map[string]bool
	virtualNodes int
	mu           sync.RWMutex
}

// NewRing creates a new Ring with the specified number of virtual nodes per physical node.
func NewRing(virtualNodes int) *Ring {
	return &Ring{
		entries:      make([]ringEntry, 0, 256),
		nodeSet:      make(map[string]bool),
		virtualNodes: virtualNodes,
	}
}

// Add places a node onto the ring by inserting virtualNodes entries.
// It is idempotent: adding the same nodeID twice is a no-op.
func (r *Ring) Add(nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.nodeSet[nodeID] {
		return // already present
	}
	r.nodeSet[nodeID] = true

	// Insert virtualNodes entries for this physical node using binary insertion.
	for i := 0; i < r.virtualNodes; i++ {
		label := fmt.Sprintf("%s#%d", nodeID, i)
		hash := r.hash(label)

		// Find insertion position using binary search.
		idx := sort.Search(len(r.entries), func(j int) bool {
			return r.entries[j].hash >= hash
		})

		// Insert at idx by extending slice and copying.
		r.entries = append(r.entries, ringEntry{})
		copy(r.entries[idx+1:], r.entries[idx:])
		r.entries[idx] = ringEntry{hash: hash, nodeID: nodeID}
	}
}

// Remove deletes all virtual nodes for nodeID from the ring.
func (r *Ring) Remove(nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.nodeSet[nodeID] {
		return // not present
	}
	delete(r.nodeSet, nodeID)

	// Filter out all entries belonging to this nodeID.
	filtered := make([]ringEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		if entry.nodeID != nodeID {
			filtered = append(filtered, entry)
		}
	}
	r.entries = filtered
}

// GetN returns up to n distinct physical node IDs starting from the successor
// of hash(key). Used to select primary + replicas.
// Returns ErrNoAliveNodes if the ring is empty.
func (r *Ring) GetN(key string, n int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.entries) == 0 {
		return nil, ErrNoAliveNodes
	}

	keyHash := r.hash(key)

	// Binary search to find the first entry >= keyHash.
	pos := sort.Search(len(r.entries), func(i int) bool {
		return r.entries[i].hash >= keyHash
	})

	if pos == len(r.entries) {
		pos = 0 // wrap around to the beginning
	}

	// Collect up to n distinct physical node IDs using fixed array (n <= 8).
	var result [8]string
	var seen [8]string
	var resultCount, seenCount int

	for i := 0; i < len(r.entries) && resultCount < n; i++ {
		idx := (pos + i) % len(r.entries)
		nodeID := r.entries[idx].nodeID

		// Check if already seen (linear search in small fixed array).
		found := false
		for j := 0; j < seenCount; j++ {
			if seen[j] == nodeID {
				found = true
				break
			}
		}

		if !found {
			seen[seenCount] = nodeID
			seenCount++
			result[resultCount] = nodeID
			resultCount++
		}
	}

	return result[:resultCount], nil
}

// Get returns the primary node for key (convenience wrapper for GetN(key, 1)).
func (r *Ring) Get(key string) (string, error) {
	nodes, err := r.GetN(key, 1)
	if err != nil {
		return "", err
	}
	if len(nodes) == 0 {
		return "", ErrNoAliveNodes
	}
	return nodes[0], nil
}

// Nodes returns a deduplicated list of all physical node IDs in the ring.
func (r *Ring) Nodes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]string, 0, len(r.nodeSet))
	for nodeID := range r.nodeSet {
		result = append(result, nodeID)
	}
	return result
}

// hash computes a deterministic uint32 from a string using FNV-1a.
// FNV-1a is 5-10x faster than SHA-256 for non-cryptographic hashing.
func (r *Ring) hash(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}
