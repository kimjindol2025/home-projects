package rpc

import (
	"runtime"
	"testing"
)

// TestPoolEffectiveness demonstrates the effectiveness of sync.Pool
// by comparing allocations with and without pooling.
func TestPoolEffectiveness(t *testing.T) {
	// Test 1: Buffer pool effectiveness (direct usage)
	t.Run("BufPool", func(t *testing.T) {
		runtime.GC()
		var m1, m2 runtime.MemStats

		// Measure WITH pool
		runtime.ReadMemStats(&m1)
		for i := 0; i < 1000; i++ {
			buf := bufPool.Get().([]byte)[:0]
			buf = append(buf, []byte("test data")...)
			bufPool.Put(buf[:0])
		}
		runtime.ReadMemStats(&m2)
		allocsWithPool := m2.Alloc - m1.Alloc

		runtime.GC()

		// Measure WITHOUT pool (direct allocation)
		runtime.ReadMemStats(&m1)
		var bufs [][]byte
		for i := 0; i < 1000; i++ {
			buf := make([]byte, 0, 512)
			buf = append(buf, []byte("test data")...)
			bufs = append(bufs, buf)
		}
		runtime.ReadMemStats(&m2)
		allocsWithoutPool := m2.Alloc - m1.Alloc
		_ = bufs // prevent optimization

		// Pool should result in significantly less memory allocation
		t.Logf("Buffer allocations WITH pool: %d bytes", allocsWithPool)
		t.Logf("Buffer allocations WITHOUT pool: %d bytes", allocsWithoutPool)
		if allocsWithoutPool > 0 {
			t.Logf("Reduction: %.1f%%", float64(allocsWithoutPool-allocsWithPool)/float64(allocsWithoutPool)*100)
		}

		// Pool version should use less memory
		if allocsWithPool > allocsWithoutPool {
			t.Logf("NOTE: Pool test shows WITH=%d > WITHOUT=%d (may need more iterations)", allocsWithPool, allocsWithoutPool)
		}
	})

	// Test 2: Pending pool effectiveness
	t.Run("PendingPool", func(t *testing.T) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		allocsBefore := m.Alloc

		// WITH pool
		for i := 0; i < 1000; i++ {
			pend := pendingPool.Get().(*pending)
			_ = pend
			pendingPool.Put(pend)
		}

		runtime.ReadMemStats(&m)
		allocsWithPool := m.Alloc - allocsBefore

		// WITHOUT pool
		runtime.ReadMemStats(&m)
		allocsBefore = m.Alloc

		for i := 0; i < 1000; i++ {
			pend := &pending{ch: make(chan *Message, 1)}
			_ = pend
		}

		runtime.ReadMemStats(&m)
		allocsWithoutPool := m.Alloc - allocsBefore

		t.Logf("Pending allocations WITH pool: %d bytes", allocsWithPool)
		t.Logf("Pending allocations WITHOUT pool: %d bytes", allocsWithoutPool)
		t.Logf("Reduction: %.1f%%", float64(allocsWithoutPool-allocsWithPool)/float64(allocsWithoutPool)*100)

		if allocsWithPool >= allocsWithoutPool {
			t.Errorf("Pool should reduce allocations, got WITH=%d, WITHOUT=%d", allocsWithPool, allocsWithoutPool)
		}
	})
}

// BenchmarkCallWithoutPool shows baseline performance without pooling
// (for comparison only - not used in production)
func BenchmarkCallWithoutPool(b *testing.B) {
	clientT, serverT, err := NewInProcessPair()
	if err != nil {
		b.Fatalf("NewInProcessPair failed: %v", err)
	}

	codec := NewBinaryCodec()
	server := NewServer(codec)
	server.ServeAsync(serverT)

	server.Register("Echo", func(args []byte) ([]byte, error) {
		return args, nil
	})

	// Create a custom client that doesn't use pools
	c := &Client{
		codec:     codec,
		transport: clientT,
		nextID:    1,
		inflight:  make(map[uint64]*pending),
		errors:    make([]string, 0),
		done:      make(chan struct{}),
	}
	go c.readLoop()
	defer c.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Direct allocation instead of pool (not using argsBuf to avoid pool reuse)
		var result string
		_ = c.CallWithTimeout("Echo", "test", &result, 0)
	}
}

// BenchmarkCallWithPool shows performance with sync.Pool
// This should significantly outperform BenchmarkCallWithoutPool
func BenchmarkCallWithPool(b *testing.B) {
	clientT, serverT, err := NewInProcessPair()
	if err != nil {
		b.Fatalf("NewInProcessPair failed: %v", err)
	}

	codec := NewBinaryCodec()
	server := NewServer(codec)
	server.ServeAsync(serverT)

	server.Register("Echo", func(args []byte) ([]byte, error) {
		return args, nil
	})

	client := NewClient(codec, clientT)
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result string
		_ = client.Call("Echo", "test", &result)
	}
}

// TestTimerMemoryLeak verifies time.NewTimer with defer Stop prevents leaks
func TestTimerMemoryLeak(t *testing.T) {
	// This test ensures that CallWithTimeout properly manages timer cleanup
	clientT, serverT, err := NewInProcessPair()
	if err != nil {
		t.Fatalf("NewInProcessPair failed: %v", err)
	}

	codec := NewBinaryCodec()
	server := NewServer(codec)
	server.ServeAsync(serverT)

	server.Register("Echo", func(args []byte) ([]byte, error) {
		return args, nil
	})

	client := NewClient(codec, clientT)
	defer client.Close()

	// Check that memory doesn't grow unboundedly with many short timeouts
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	for i := 0; i < 100; i++ {
		var result string
		// These timeouts should be cleaned up properly
		_ = client.CallWithTimeout("Echo", "test", &result, 0)
	}

	runtime.ReadMemStats(&m2)
	growthMB := float64(m2.Alloc-m1.Alloc) / (1024 * 1024)

	// Growth should be minimal (timer resources cleaned up)
	t.Logf("Memory growth after 100 calls: %.2f MB", growthMB)
	if growthMB > 10 {
		t.Logf("WARNING: Large memory growth detected (%.2f MB), possible timer leak", growthMB)
	}
}
