package kvstore

import "testing"

// Helper: create a test cluster with n nodes.
func newTestCluster(n int, rf int) *Cluster {
	c := NewCluster(rf)
	for i := 0; i < n; i++ {
		id := string(rune('A') + rune(i))
		addr := id + ":7001"
		c.AddNode(id, addr)
	}
	return c
}

// ============ Ring Tests ============

func TestRingAdd(t *testing.T) {
	r := NewRing(10)
	r.Add("node-1")

	nodes := r.Nodes()
	if len(nodes) != 1 || nodes[0] != "node-1" {
		t.Errorf("Expected [node-1], got %v", nodes)
	}
}

func TestRingAddIdempotent(t *testing.T) {
	r := NewRing(10)
	r.Add("node-1")
	r.Add("node-1") // second add should be no-op

	nodes := r.Nodes()
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}
}

func TestRingRemove(t *testing.T) {
	r := NewRing(10)
	r.Add("node-1")
	r.Add("node-2")

	r.Remove("node-1")

	nodes := r.Nodes()
	if len(nodes) != 1 || nodes[0] != "node-2" {
		t.Errorf("Expected [node-2], got %v", nodes)
	}
}

func TestRingGetPrimary(t *testing.T) {
	r := NewRing(10)
	r.Add("node-1")
	r.Add("node-2")
	r.Add("node-3")

	primary, err := r.Get("somekey")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Primary should be one of the nodes.
	valid := false
	for _, node := range r.Nodes() {
		if primary == node {
			valid = true
			break
		}
	}
	if !valid {
		t.Errorf("Primary %s not in node set", primary)
	}
}

func TestRingGetN(t *testing.T) {
	r := NewRing(10)
	r.Add("A")
	r.Add("B")
	r.Add("C")
	r.Add("D")

	replicas, err := r.GetN("testkey", 3)
	if err != nil {
		t.Fatalf("GetN failed: %v", err)
	}

	if len(replicas) != 3 {
		t.Errorf("Expected 3 replicas, got %d: %v", len(replicas), replicas)
	}

	// Check all are distinct.
	seen := make(map[string]bool)
	for _, replica := range replicas {
		if seen[replica] {
			t.Errorf("Duplicate replica: %s", replica)
		}
		seen[replica] = true
	}
}

func TestRingGetNFewer(t *testing.T) {
	r := NewRing(10)
	r.Add("A")
	r.Add("B")

	replicas, err := r.GetN("testkey", 5)
	if err != nil {
		t.Fatalf("GetN failed: %v", err)
	}

	if len(replicas) != 2 {
		t.Errorf("Expected 2 replicas (only 2 nodes exist), got %d: %v", len(replicas), replicas)
	}
}

func TestRingDistribution(t *testing.T) {
	r := NewRing(150)
	r.Add("A")
	r.Add("B")
	r.Add("C")
	r.Add("D")
	r.Add("E")

	counts := make(map[string]int)
	for i := 0; i < 10000; i++ {
		key := "key-" + string(rune(i))
		primary, _ := r.Get(key)
		counts[primary]++
	}

	// With 5 nodes and 10000 keys, expect each to have ~2000 keys.
	// Allow 30% variance.
	expectedPerNode := 10000 / 5
	maxVariance := int(float64(expectedPerNode) * 0.3)

	for node, count := range counts {
		if count < expectedPerNode-maxVariance || count > expectedPerNode+maxVariance {
			t.Errorf("Node %s has %d keys (expected ~%d, variance %d)", node, count, expectedPerNode, maxVariance)
		}
	}
}

func TestRingEmpty(t *testing.T) {
	r := NewRing(10)

	_, err := r.Get("key")
	if err != ErrNoAliveNodes {
		t.Errorf("Expected ErrNoAliveNodes, got %v", err)
	}
}

// ============ Store Tests ============

func TestStoreSetGet(t *testing.T) {
	s := newStore()

	s.set("key1", []byte("value1"))

	val, ok := s.get("key1")
	if !ok || string(val) != "value1" {
		t.Errorf("Expected value1, got %v (ok=%v)", val, ok)
	}
}

func TestStoreGetMissing(t *testing.T) {
	s := newStore()

	val, ok := s.get("missing")
	if ok || val != nil {
		t.Errorf("Expected (nil, false), got (%v, %v)", val, ok)
	}
}

func TestStoreDelete(t *testing.T) {
	s := newStore()

	s.set("key1", []byte("value1"))
	ok := s.delete("key1")
	if !ok {
		t.Errorf("Delete should have returned true")
	}

	val, ok := s.get("key1")
	if ok || val != nil {
		t.Errorf("After delete, expected (nil, false), got (%v, %v)", val, ok)
	}
}

func TestStoreSize(t *testing.T) {
	s := newStore()

	if s.size() != 0 {
		t.Errorf("Expected size 0, got %d", s.size())
	}

	s.set("key1", []byte("value1"))
	if s.size() != 1 {
		t.Errorf("Expected size 1, got %d", s.size())
	}

	s.set("key2", []byte("value2"))
	if s.size() != 2 {
		t.Errorf("Expected size 2, got %d", s.size())
	}

	s.delete("key1")
	if s.size() != 1 {
		t.Errorf("Expected size 1 after delete, got %d", s.size())
	}
}

// ============ Cluster Tests ============

func TestClusterAddNode(t *testing.T) {
	c := NewCluster(3)
	err := c.AddNode("A", "A:7001")
	if err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}

	stats := c.Stats()
	if stats.TotalNodes != 1 {
		t.Errorf("Expected 1 node, got %d", stats.TotalNodes)
	}

	status, _ := c.NodeStatus("A")
	if status != NodeAlive {
		t.Errorf("Expected NodeAlive, got %v", status)
	}
}

func TestClusterAddNodeDuplicate(t *testing.T) {
	c := NewCluster(3)
	c.AddNode("A", "A:7001")
	err := c.AddNode("A", "A:7001")
	if err == nil {
		t.Errorf("Expected error adding duplicate node, got nil")
	}
}

func TestClusterRemoveNode(t *testing.T) {
	c := newTestCluster(3, 3)

	err := c.RemoveNode("A")
	if err != nil {
		t.Fatalf("RemoveNode failed: %v", err)
	}

	stats := c.Stats()
	if stats.TotalNodes != 2 {
		t.Errorf("Expected 2 nodes after remove, got %d", stats.TotalNodes)
	}

	_, err = c.NodeStatus("A")
	if err != ErrNodeNotFound {
		t.Errorf("Expected ErrNodeNotFound, got %v", err)
	}
}

func TestClusterSetGet(t *testing.T) {
	c := newTestCluster(3, 3)

	err := c.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	val, err := c.Get("key1")
	if err != nil || val != "value1" {
		t.Errorf("Expected value1, got %v (err=%v)", val, err)
	}
}

func TestClusterGetMissing(t *testing.T) {
	c := newTestCluster(3, 3)

	_, err := c.Get("missing")
	if err != ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestClusterDeleteKey(t *testing.T) {
	c := newTestCluster(3, 3)

	c.Set("key1", "value1")

	err := c.Delete("key1")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = c.Get("key1")
	if err != ErrKeyNotFound {
		t.Errorf("After delete, expected ErrKeyNotFound, got %v", err)
	}
}

func TestClusterDeleteMissing(t *testing.T) {
	c := newTestCluster(3, 3)

	err := c.Delete("missing")
	if err != ErrKeyNotFound {
		t.Errorf("Delete on missing key should return ErrKeyNotFound, got %v", err)
	}
}

func TestClusterReplicationFactor3(t *testing.T) {
	c := newTestCluster(5, 3)

	err := c.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Check that key exists in exactly 3 stores (or fewer if some nodes don't have it).
	// For a consistent hash, the key should be in its 3 replicas.
	replicas, _ := c.ring.GetN("key1", 3)
	foundCount := 0
	for _, nodeID := range replicas {
		store := c.stores[nodeID]
		if _, ok := store.get("key1"); ok {
			foundCount++
		}
	}

	if foundCount != 3 {
		t.Errorf("Expected key in 3 replicas, found in %d", foundCount)
	}
}

func TestClusterReplicationOnDeadNode(t *testing.T) {
	c := newTestCluster(3, 3)

	// Mark one node as dead.
	c.SimulateHeartbeat("A", false)

	// Write should still succeed because 2/3 alive nodes = quorum.
	err := c.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Set with 2/3 alive nodes failed: %v", err)
	}

	// Read should still work.
	val, err := c.Get("key1")
	if err != nil || val != "value1" {
		t.Errorf("Get failed: %v", err)
	}
}

func TestClusterReadFallback(t *testing.T) {
	c := newTestCluster(3, 3)

	c.Set("key1", "value1")

	// Mark primary as dead.
	primary, _ := c.ring.Get("key1")
	c.SimulateHeartbeat(primary, false)

	// Read should still succeed from a replica.
	val, err := c.Get("key1")
	if err != nil || val != "value1" {
		t.Errorf("Read fallback failed: %v", err)
	}
}

func TestClusterReadFailsAllDead(t *testing.T) {
	c := newTestCluster(3, 3)

	c.Set("key1", "value1")

	// Mark all nodes as dead.
	c.SimulateHeartbeat("A", false)
	c.SimulateHeartbeat("B", false)
	c.SimulateHeartbeat("C", false)

	// Read should fail because no alive replicas.
	_, err := c.Get("key1")
	if err != ErrKeyNotFound && err != ErrNoAliveNodes {
		t.Errorf("Expected ErrKeyNotFound or ErrNoAliveNodes, got %v", err)
	}
}

func TestClusterHeartbeatMarksDead(t *testing.T) {
	c := newTestCluster(1, 1)

	// Initially alive.
	status, _ := c.NodeStatus("A")
	if status != NodeAlive {
		t.Errorf("Expected NodeAlive initially, got %v", status)
	}

	// Simulate missed heartbeats (directly set node to dead).
	c.SimulateHeartbeat("A", false)

	status, _ = c.NodeStatus("A")
	if status != NodeDead {
		t.Errorf("Expected NodeDead after marking dead, got %v", status)
	}
}

func TestClusterHeartbeatRecovery(t *testing.T) {
	c := newTestCluster(1, 1)

	// Mark as dead.
	c.SimulateHeartbeat("A", false)

	status, _ := c.NodeStatus("A")
	if status != NodeDead {
		t.Errorf("Expected NodeDead, got %v", status)
	}

	// Record a successful heartbeat to recover.
	c.SimulateHeartbeat("A", true)

	status, _ = c.NodeStatus("A")
	if status != NodeAlive {
		t.Errorf("Expected NodeAlive after recovery, got %v", status)
	}
}

func TestClusterDataMigrationOnAdd(t *testing.T) {
	c := NewCluster(3)
	c.AddNode("A", "A:7001")
	c.AddNode("B", "B:7001")

	// Write some data.
	c.Set("key1", "value1")
	c.Set("key2", "value2")

	// Add a third node.
	c.AddNode("C", "C:7001")

	// All keys should still be readable (possibly from new node).
	val1, err := c.Get("key1")
	if err != nil || val1 != "value1" {
		t.Errorf("Get key1 after migration failed: %v", err)
	}

	val2, err := c.Get("key2")
	if err != nil || val2 != "value2" {
		t.Errorf("Get key2 after migration failed: %v", err)
	}
}

func TestClusterDataMigrationOnRemove(t *testing.T) {
	c := newTestCluster(3, 3)

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	// Remove a node.
	c.RemoveNode("A")

	// All keys should still be readable (from remaining replicas).
	val1, err := c.Get("key1")
	if err != nil || val1 != "value1" {
		t.Errorf("Get key1 after removal failed: %v", err)
	}

	val2, err := c.Get("key2")
	if err != nil || val2 != "value2" {
		t.Errorf("Get key2 after removal failed: %v", err)
	}
}

func TestClusterEmptyKey(t *testing.T) {
	c := newTestCluster(3, 3)

	err := c.Set("", "value")
	if err != ErrEmptyKey {
		t.Errorf("Expected ErrEmptyKey, got %v", err)
	}
}

func TestClusterStats(t *testing.T) {
	c := newTestCluster(3, 3)

	stats := c.Stats()
	if stats.TotalNodes != 3 {
		t.Errorf("Expected 3 nodes, got %d", stats.TotalNodes)
	}
	if stats.AliveNodes != 3 {
		t.Errorf("Expected 3 alive nodes, got %d", stats.AliveNodes)
	}

	c.Set("key1", "value1")
	stats = c.Stats()
	// Total keys counts replicas, so 3 replicas = 3 keys.
	if stats.TotalKeys != 3 {
		t.Errorf("Expected 3 total keys (across replicas), got %d", stats.TotalKeys)
	}
}

func TestClusterNodeStatusAfterDead(t *testing.T) {
	c := newTestCluster(1, 1)

	c.SimulateHeartbeat("A", false)

	status, _ := c.NodeStatus("A")
	if status != NodeDead {
		t.Errorf("Expected NodeDead, got %v", status)
	}
}
