package rpc

import (
	"sync"
	"testing"
)

// BenchmarkClientCall_serial benchmarks serial RPC calls.
func BenchmarkClientCall_serial(b *testing.B) {
	clientT, serverT, err := NewInProcessPair()
	if err != nil {
		b.Fatalf("NewInProcessPair failed: %v", err)
	}

	codec := NewBinaryCodec()
	server := NewServer(codec)
	server.ServeAsync(serverT)

	// Register simple echo handler
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

// BenchmarkClientCall_parallel benchmarks parallel RPC calls (10 goroutines).
func BenchmarkClientCall_parallel(b *testing.B) {
	clientT, serverT, err := NewInProcessPair()
	if err != nil {
		b.Fatalf("NewInProcessPair failed: %v", err)
	}

	codec := NewBinaryCodec()
	server := NewServer(codec)
	server.ServeAsync(serverT)

	// Register simple echo handler
	server.Register("Echo", func(args []byte) ([]byte, error) {
		return args, nil
	})

	client := NewClient(codec, clientT)
	defer client.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var result string
			_ = client.Call("Echo", "test", &result)
		}
	})
}

// BenchmarkCodecMarshal_string benchmarks string marshaling.
func BenchmarkCodecMarshal_string(b *testing.B) {
	codec := NewBinaryCodec()
	testData := "hello world this is a test string"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = codec.Marshal(testData, make([]byte, 0, 64))
	}
}

// BenchmarkCodecMarshal_bytes benchmarks []byte marshaling.
func BenchmarkCodecMarshal_bytes(b *testing.B) {
	codec := NewBinaryCodec()
	testData := []byte("hello world this is a test string")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = codec.Marshal(testData, make([]byte, 0, 64))
	}
}

// BenchmarkCodecUnmarshal_string benchmarks string unmarshaling.
func BenchmarkCodecUnmarshal_string(b *testing.B) {
	codec := NewBinaryCodec()
	testData := "hello world"
	encoded, _ := codec.Marshal(testData, make([]byte, 0, 64))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result string
		_ = codec.Unmarshal(encoded, &result)
	}
}

// BenchmarkBufferPoolReuse benchmarks buffer pool reuse (sync.Pool efficiency).
func BenchmarkBufferPoolReuse(b *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 512)
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := pool.Get().([]byte)[:0]
		_ = append(buf, []byte("test data")...)
		pool.Put(buf[:0])
	}
}

// Note: NewRing is in pkg/kvstore, not pkg/rpc,
// so we don't benchmark it here to avoid cross-package dependencies.
