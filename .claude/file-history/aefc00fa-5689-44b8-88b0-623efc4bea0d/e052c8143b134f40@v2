package dispatcher

import (
	"testing"
)

func TestRingBuffer_Create(t *testing.T) {
	rb, err := NewRingBuffer(1024)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if rb.Capacity() != 1024 {
		t.Fatalf("Capacity = %d, want 1024", rb.Capacity())
	}

	if !rb.IsEmpty() {
		t.Fatal("Buffer should be empty after creation")
	}
}

func TestRingBuffer_PowerOf2Check(t *testing.T) {
	_, err := NewRingBuffer(1000) // Not power of 2
	if err == nil {
		t.Fatal("Should reject non-power-of-2 capacity")
	}
}

func TestRingBuffer_Publish_And_Subscribe(t *testing.T) {
	rb, _ := NewRingBuffer(16)

	job := Job{
		ID:       1,
		Priority: 5,
		Payload:  []byte("test payload"),
	}

	// Publish
	if !rb.Publish(job) {
		t.Fatal("Publish should succeed")
	}

	if rb.IsEmpty() {
		t.Fatal("Buffer should not be empty after publish")
	}

	// Subscribe
	result, ok := rb.Subscribe()
	if !ok {
		t.Fatal("Subscribe should return data")
	}

	if result.ID != job.ID || string(result.Payload) != string(job.Payload) {
		t.Fatal("Published and subscribed data don't match")
	}

	if !rb.IsEmpty() {
		t.Fatal("Buffer should be empty after subscribe")
	}
}

func TestRingBuffer_Size(t *testing.T) {
	rb, _ := NewRingBuffer(16)

	if rb.Size() != 0 {
		t.Fatalf("Size = %d, want 0", rb.Size())
	}

	rb.Publish(Job{ID: 1})
	rb.Publish(Job{ID: 2})
	rb.Publish(Job{ID: 3})

	if rb.Size() != 3 {
		t.Fatalf("Size = %d, want 3", rb.Size())
	}

	rb.Subscribe()
	if rb.Size() != 2 {
		t.Fatalf("Size after subscribe = %d, want 2", rb.Size())
	}
}

func TestRingBuffer_Full(t *testing.T) {
	rb, _ := NewRingBuffer(4)

	// Fill the buffer
	for i := 0; i < 4; i++ {
		if !rb.Publish(Job{ID: uint64(i)}) {
			t.Fatalf("Publish %d should succeed", i)
		}
	}

	if !rb.IsFull() {
		t.Fatal("Buffer should be full")
	}

	// Try to overflow
	if rb.Publish(Job{ID: 999}) {
		t.Fatal("Publish to full buffer should fail")
	}
}

func TestRingBuffer_Wrap(t *testing.T) {
	rb, _ := NewRingBuffer(8)

	// Publish and subscribe several times to test wrapping
	for round := 0; round < 3; round++ {
		for i := 0; i < 5; i++ {
			if !rb.Publish(Job{ID: uint64(round*5 + i)}) {
				t.Fatalf("Publish round %d, job %d failed", round, i)
			}
		}

		for i := 0; i < 5; i++ {
			job, ok := rb.Subscribe()
			if !ok {
				t.Fatalf("Subscribe round %d, job %d failed", round, i)
			}
			if job.ID != uint64(round*5+i) {
				t.Fatalf("Job ID mismatch: got %d, want %d", job.ID, round*5+i)
			}
		}
	}
}

func TestRingBuffer_SubscribeBatch(t *testing.T) {
	rb, _ := NewRingBuffer(64)

	// Publish 10 items
	for i := 0; i < 10; i++ {
		rb.Publish(Job{ID: uint64(i), Payload: []byte{byte(i)}})
	}

	// Subscribe batch of 5
	batch := rb.SubscribeBatch(5)
	if len(batch) != 5 {
		t.Fatalf("Batch size = %d, want 5", len(batch))
	}

	for i := 0; i < 5; i++ {
		if batch[i].ID != uint64(i) {
			t.Fatalf("Batch[%d].ID = %d, want %d", i, batch[i].ID, i)
		}
	}

	// Subscribe remaining 5
	batch2 := rb.SubscribeBatch(10)
	if len(batch2) != 5 {
		t.Fatalf("Second batch size = %d, want 5", len(batch2))
	}

	if !rb.IsEmpty() {
		t.Fatal("Buffer should be empty")
	}
}

func TestRingBuffer_Reset(t *testing.T) {
	rb, _ := NewRingBuffer(16)

	rb.Publish(Job{ID: 1})
	rb.Publish(Job{ID: 2})

	if rb.IsEmpty() {
		t.Fatal("Buffer should have items")
	}

	rb.Reset()

	if !rb.IsEmpty() {
		t.Fatal("Buffer should be empty after reset")
	}

	if rb.Size() != 0 {
		t.Fatalf("Size after reset = %d, want 0", rb.Size())
	}
}

func BenchmarkRingBuffer_Publish(b *testing.B) {
	rb, _ := NewRingBuffer(65536)
	job := Job{ID: 1, Payload: []byte("test")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Publish(job)
		if rb.IsFull() {
			rb.Reset()
		}
	}
}

func BenchmarkRingBuffer_Subscribe(b *testing.B) {
	rb, _ := NewRingBuffer(65536)

	// Pre-fill
	for i := 0; i < 1000; i++ {
		rb.Publish(Job{ID: uint64(i)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Subscribe()
		if rb.IsEmpty() {
			for j := 0; j < 1000; j++ {
				rb.Publish(Job{ID: uint64(j)})
			}
		}
	}
}

func BenchmarkRingBuffer_PublishSubscribe(b *testing.B) {
	rb, _ := NewRingBuffer(65536)
	job := Job{ID: 1, Payload: []byte("test data")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Publish(job)
		rb.Subscribe()
	}
}
