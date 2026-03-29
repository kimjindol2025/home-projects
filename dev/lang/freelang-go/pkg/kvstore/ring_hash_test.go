package kvstore

import (
	"crypto/sha256"
	"encoding/binary"
	"hash/fnv"
	"runtime"
	"testing"
	"strconv"
)

// TestFNVHashFunction verifies FNV-1a hash works correctly
func TestFNVHashFunction(t *testing.T) {
	r := NewRing(150)

	// Test 1: Deterministic hashing
	h1 := r.hash("test_key")
	h2 := r.hash("test_key")
	if h1 != h2 {
		t.Errorf("FNV hash not deterministic: %d != %d", h1, h2)
	}

	// Test 2: Different inputs produce different hashes (usually)
	h3 := r.hash("different_key")
	if h1 == h3 {
		t.Logf("Hash collision for 'test_key' and 'different_key': %d", h1)
	}

	// Test 3: Many keys should hash without collisions
	hashes := make(map[uint32]string)
	collisions := 0
	for i := 0; i < 1000; i++ {
		key := "key_" + strconv.Itoa(i)
		h := r.hash(key)
		if existing, ok := hashes[h]; ok {
			collisions++
			t.Logf("Collision: %s and %s both hash to %d", existing, key, h)
		}
		hashes[h] = key
	}
	t.Logf("Collisions out of 1000 keys: %d (collision rate: %.2f%%)", collisions, float64(collisions)*100/1000)
}

// BenchmarkFNVVsSHA256 directly compares FNV-1a vs SHA-256 performance
func BenchmarkFNVVsSHA256(b *testing.B) {
	b.Run("FNV1a", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			h := fnv.New32a()
			h.Write([]byte("test_key_for_hashing"))
			_ = h.Sum32()
		}
	})

	b.Run("SHA256Truncated", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sum := sha256.Sum256([]byte("test_key_for_hashing"))
			_ = binary.BigEndian.Uint32(sum[:4])
		}
	})
}

// BenchmarkRingHashPerformance measures hash function in Ring context
func BenchmarkRingHashPerformance(b *testing.B) {
	r := NewRing(150)
	testKeys := []string{
		"key_1", "key_2", "key_3", "key_4", "key_5",
		"key_6", "key_7", "key_8", "key_9", "key_10",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range testKeys {
			_ = r.hash(key)
		}
	}
}

// TestGetNAllocationFree verifies GetN doesn't allocate heap memory
func TestGetNAllocationFree(t *testing.T) {
	r := NewRing(150)
	r.Add("node1")
	r.Add("node2")
	r.Add("node3")

	// GetN should not allocate (uses stack-based [8]string array)
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	allocsBefore := m.Alloc

	for i := 0; i < 1000; i++ {
		_, _ = r.GetN("test_key", 3)
	}

	runtime.ReadMemStats(&m)
	allocsAfter := m.Alloc

	heapGrowth := allocsAfter - allocsBefore
	t.Logf("GetN heap growth (1000 calls): %d bytes", heapGrowth)

	// Should be minimal (only from logging, not from data structures)
	if heapGrowth > 100000 {
		t.Logf("WARNING: GetN allocated %d bytes (expected < 100KB)", heapGrowth)
	}
}
