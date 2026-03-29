package kvstore

import (
	"fmt"
	"testing"
)

// BenchmarkRingHash benchmarks the hash function performance.
func BenchmarkRingHash(b *testing.B) {
	r := NewRing(150)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = r.hash(fmt.Sprintf("key_%d", i))
	}
}

// BenchmarkRingGetN benchmarks GetN operation (consistent hashing lookup).
func BenchmarkRingGetN(b *testing.B) {
	r := NewRing(150)

	// Setup: add 3 nodes
	r.Add("node1")
	r.Add("node2")
	r.Add("node3")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.GetN(fmt.Sprintf("key_%d", i), 3)
	}
}

// BenchmarkClusterSet benchmarks Set operation.
func BenchmarkClusterSet(b *testing.B) {
	cluster := NewCluster(3)
	cluster.AddNode("A", "A:7001")
	cluster.AddNode("B", "B:7001")
	cluster.AddNode("C", "C:7001")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		val := fmt.Sprintf("value_%d", i)
		_ = cluster.Set(key, val)
	}
}

// BenchmarkClusterGet benchmarks Get operation.
func BenchmarkClusterGet(b *testing.B) {
	cluster := NewCluster(3)
	cluster.AddNode("A", "A:7001")
	cluster.AddNode("B", "B:7001")
	cluster.AddNode("C", "C:7001")

	// Populate with data
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		val := fmt.Sprintf("value_%d", i)
		cluster.Set(key, val)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		_, _ = cluster.Get(key)
	}
}
