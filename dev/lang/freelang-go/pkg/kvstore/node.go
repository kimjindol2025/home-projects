package kvstore

import (
	"sync"
	"time"
)

// NodeStatus represents the health state of a node.
type NodeStatus int

const (
	NodeAlive NodeStatus = iota
	NodeSuspect
	NodeDead
)

// Node represents a single storage peer in the cluster.
type Node struct {
	ID            string
	Addr          string
	Status        NodeStatus
	lastHeartbeat time.Time
	mu            sync.RWMutex
}

// NewNode creates a new node with the given ID and address.
func NewNode(id, addr string) *Node {
	return &Node{
		ID:            id,
		Addr:          addr,
		Status:        NodeAlive,
		lastHeartbeat: time.Now(),
	}
}

// MarkAlive transitions the node to the Alive state.
func (n *Node) MarkAlive() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Status = NodeAlive
}

// MarkSuspect transitions the node to the Suspect state.
func (n *Node) MarkSuspect() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Status = NodeSuspect
}

// MarkDead transitions the node to the Dead state.
func (n *Node) MarkDead() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Status = NodeDead
}

// IsAlive returns true if the node is in Alive state.
func (n *Node) IsAlive() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Status == NodeAlive
}

// HeartbeatAge returns the duration since the last successful heartbeat.
func (n *Node) HeartbeatAge() time.Duration {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return time.Since(n.lastHeartbeat)
}

// RecordHeartbeat updates the lastHeartbeat timestamp to now.
func (n *Node) RecordHeartbeat() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.lastHeartbeat = time.Now()
}
