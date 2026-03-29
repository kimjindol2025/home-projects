package kvstore

import (
	"fmt"
	"sync"
	"time"
)

const (
	ReplicationFactor        = 3
	DefaultHeartbeatInterval = 5 * time.Second
	HeartbeatMissThreshold   = 3
)

// Cluster orchestrates a set of Nodes, a Ring, and per-node stores.
// It is the only type callers interact with.
type Cluster struct {
	nodes   map[string]*Node
	stores  map[string]*store
	ring    *Ring
	rf      int
	mu      sync.RWMutex
	stopCh  chan struct{}
	errors  []string
	ticker  *time.Ticker
	done    chan struct{}
}

// ClusterStats provides a snapshot of cluster health.
type ClusterStats struct {
	TotalNodes int
	AliveNodes int
	DeadNodes  int
	TotalKeys  int
}

// NewCluster creates a Cluster with the given replication factor.
func NewCluster(rf int) *Cluster {
	return &Cluster{
		nodes:   make(map[string]*Node),
		stores:  make(map[string]*store),
		ring:    NewRing(DefaultVirtualNodes),
		rf:      rf,
		stopCh:  make(chan struct{}),
		errors:  make([]string, 0),
		done:    make(chan struct{}),
	}
}

// AddNode registers a new node and triggers ring rebalancing + data migration.
func (c *Cluster) AddNode(id, addr string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.nodes[id]; exists {
		return fmt.Errorf("node %s already exists", id)
	}

	// Create node and store.
	node := NewNode(id, addr)
	store := newStore()

	c.nodes[id] = node
	c.stores[id] = store

	// Add to ring.
	c.ring.Add(id)

	// Migrate data: for each key in all other stores, check if it should live on the new node.
	for nodeID, s := range c.stores {
		if nodeID == id {
			continue
		}
		keys := s.keys()
		for _, key := range keys {
			// Check if key now belongs to the new node.
			replicas, _ := c.ring.GetN(key, c.rf)
			shouldMigrate := false
			for _, replica := range replicas {
				if replica == id {
					shouldMigrate = true
					break
				}
			}
			if shouldMigrate {
				val, _ := s.get(key)
				c.stores[id].set(key, val)
			}
		}
	}

	return nil
}

// RemoveNode deregisters a node, triggering data migration to surviving replicas.
func (c *Cluster) RemoveNode(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.nodes[id]; !exists {
		return ErrNodeNotFound
	}

	// Migrate data from the removed node to surviving replicas.
	if store, ok := c.stores[id]; ok {
		keys := store.keys()
		for _, key := range keys {
			val, _ := store.get(key)
			// Route to new replicas (excluding the removed node).
			replicas, _ := c.ring.GetN(key, c.rf)
			for _, replica := range replicas {
				if replica != id {
					if s, ok := c.stores[replica]; ok {
						s.set(key, val)
					}
				}
			}
		}
	}

	// Remove from tracking structures.
	delete(c.nodes, id)
	delete(c.stores, id)
	c.ring.Remove(id)

	return nil
}

// Set writes value to the primary node + (rf-1) replicas.
// Returns ErrNoAliveNodes if fewer than 1 replica is reachable.
// Returns ErrReplicaFailure if fewer than quorum (rf/2+1) acknowledge.
func (c *Cluster) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get replica nodes.
	replicas, err := c.ring.GetN(key, c.rf)
	if err != nil {
		return err
	}

	// Filter to alive nodes.
	alive := make([]string, 0, len(replicas))
	for _, nodeID := range replicas {
		if node, ok := c.nodes[nodeID]; ok && node.IsAlive() {
			alive = append(alive, nodeID)
		}
	}

	if len(alive) == 0 {
		return ErrNoAliveNodes
	}

	// Compute quorum requirement.
	quorum := c.rf/2 + 1

	// Write to all alive replicas.
	acks := 0
	for _, nodeID := range alive {
		if store, ok := c.stores[nodeID]; ok {
			store.set(key, []byte(value))
			acks++
		}
	}

	if acks < quorum {
		return ErrReplicaFailure
	}

	return nil
}

// Get reads from the primary node. If the primary is dead,
// falls back to the next live replica.
func (c *Cluster) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get replica nodes in order.
	replicas, err := c.ring.GetN(key, c.rf)
	if err != nil {
		return "", err
	}

	// Try each replica in order until we find an alive one with the key.
	for _, nodeID := range replicas {
		node, ok := c.nodes[nodeID]
		if !ok || !node.IsAlive() {
			continue
		}

		store, ok := c.stores[nodeID]
		if !ok {
			continue
		}

		val, found := store.get(key)
		if found {
			return string(val), nil
		}
	}

	return "", ErrKeyNotFound
}

// Delete removes a key from all replicas.
func (c *Cluster) Delete(key string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get replica nodes.
	replicas, err := c.ring.GetN(key, c.rf)
	if err != nil {
		return err
	}

	// Delete from all replicas.
	found := false
	for _, nodeID := range replicas {
		if store, ok := c.stores[nodeID]; ok {
			if store.delete(key) {
				found = true
			}
		}
	}

	if !found {
		return ErrKeyNotFound
	}

	return nil
}

// NodeStatus returns the current status of a node.
func (c *Cluster) NodeStatus(nodeID string) (NodeStatus, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	node, ok := c.nodes[nodeID]
	if !ok {
		return NodeAlive, ErrNodeNotFound
	}

	return node.Status, nil
}

// SimulateHeartbeat is used in tests to manually trigger heartbeat processing.
func (c *Cluster) SimulateHeartbeat(nodeID string, alive bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.nodes[nodeID]
	if !ok {
		return
	}

	if alive {
		node.RecordHeartbeat()
		node.MarkAlive()
	} else {
		// Simulate missed heartbeats by setting lastHeartbeat far in the past.
		node.mu.Lock()
		node.lastHeartbeat = time.Now().Add(-10 * DefaultHeartbeatInterval)
		node.mu.Unlock()
		node.MarkDead()
	}
}

// StartHeartbeat begins background goroutine for periodic heartbeat monitoring.
func (c *Cluster) StartHeartbeat(interval time.Duration) {
	c.ticker = time.NewTicker(interval)
	go func() {
		defer close(c.done)
		for {
			select {
			case <-c.ticker.C:
				c.processHeartbeat(interval)
			case <-c.stopCh:
				c.ticker.Stop()
				return
			}
		}
	}()
}

// processHeartbeat checks the health of all nodes and transitions states as needed.
func (c *Cluster) processHeartbeat(interval time.Duration) {
	c.mu.RLock()
	nodes := make([]*Node, 0, len(c.nodes))
	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}
	c.mu.RUnlock()

	for _, node := range nodes {
		age := node.HeartbeatAge()
		beatThreshold := time.Duration(HeartbeatMissThreshold) * interval

		if age < interval {
			// Recent heartbeat, ensure Alive.
			if !node.IsAlive() {
				node.MarkAlive()
			}
		} else if age < beatThreshold {
			// Missed one but not all, mark Suspect.
			if node.Status == NodeAlive {
				node.MarkSuspect()
			}
		} else {
			// Missed too many, mark Dead.
			node.mu.RLock()
			currentStatus := node.Status
			node.mu.RUnlock()

			if currentStatus != NodeDead {
				node.MarkDead()
			}
		}
	}
}

// StopHeartbeat shuts down the heartbeat goroutine gracefully.
func (c *Cluster) StopHeartbeat() {
	if c.ticker != nil {
		close(c.stopCh)
		<-c.done
	}
}

// Errors returns non-fatal accumulated errors.
func (c *Cluster) Errors() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]string(nil), c.errors...)
}

// Stats returns a snapshot of cluster health.
func (c *Cluster) Stats() ClusterStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	alive := 0
	dead := 0
	for _, node := range c.nodes {
		if node.IsAlive() {
			alive++
		} else {
			dead++
		}
	}

	totalKeys := 0
	for _, s := range c.stores {
		totalKeys += s.size()
	}

	return ClusterStats{
		TotalNodes: len(c.nodes),
		AliveNodes: alive,
		DeadNodes:  dead,
		TotalKeys:  totalKeys,
	}
}
