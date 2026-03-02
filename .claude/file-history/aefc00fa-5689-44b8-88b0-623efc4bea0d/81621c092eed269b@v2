package dispatcher

import (
	"fmt"
	"sync/atomic"
)

// Job represents a unit of work
type Job struct {
	ID       uint64
	Priority int
	Payload  []byte
	Response chan []byte // 응답 채널
}

// RingBuffer: Lock-Free SPSC (Single-Producer, Single-Consumer) Ring Buffer
// Inspired by LMAX Disruptor pattern
type RingBuffer struct {
	data     []Job
	capacity uint64
	mask     uint64 // capacity - 1 (빠른 모듈로 연산)

	// 발행자 (publisher) - Dispatcher에서만 쓰기
	publisherCursor uint64 // atomic

	// 구독자 (subscriber) - Worker에서만 읽기
	subscriberCursor uint64 // atomic

	// 패딩: False sharing 방지
	_pad1 [56]byte // Cache line padding
}

// NewRingBuffer creates a new ring buffer with capacity (must be power of 2)
func NewRingBuffer(capacity uint64) (*RingBuffer, error) {
	// Validate capacity is power of 2
	if capacity == 0 || (capacity&(capacity-1)) != 0 {
		return nil, fmt.Errorf("capacity must be power of 2, got %d", capacity)
	}

	rb := &RingBuffer{
		data:             make([]Job, capacity),
		capacity:         capacity,
		mask:             capacity - 1,
		publisherCursor:  0,
		subscriberCursor: 0,
	}

	return rb, nil
}

// Size returns current number of items in buffer
func (rb *RingBuffer) Size() uint64 {
	pub := atomic.LoadUint64(&rb.publisherCursor)
	sub := atomic.LoadUint64(&rb.subscriberCursor)
	return pub - sub
}

// IsFull checks if buffer is full
func (rb *RingBuffer) IsFull() bool {
	return rb.Size() >= rb.capacity
}

// IsEmpty checks if buffer is empty
func (rb *RingBuffer) IsEmpty() bool {
	return rb.Size() == 0
}

// Publish adds a job to the ring buffer (Publisher only)
// Returns false if buffer is full
func (rb *RingBuffer) Publish(job Job) bool {
	pub := atomic.LoadUint64(&rb.publisherCursor)
	sub := atomic.LoadUint64(&rb.subscriberCursor)

	// Check if next position would wrap into subscriber
	nextPub := pub + 1
	if (nextPub-sub) > rb.capacity {
		return false // Buffer full
	}

	// Write job to ring buffer (no synchronization needed - SPSC)
	rb.data[pub&rb.mask] = job

	// Publish cursor atomically
	atomic.StoreUint64(&rb.publisherCursor, nextPub)

	return true
}

// TryPublish attempts to publish with timeout semantics
func (rb *RingBuffer) TryPublish(job Job, maxRetries int) bool {
	for i := 0; i < maxRetries; i++ {
		if rb.Publish(job) {
			return true
		}
		// Spin or backoff (simple spin here)
	}
	return false
}

// Subscribe retrieves the next job (Subscriber only)
// Returns (nil, false) if buffer is empty
func (rb *RingBuffer) Subscribe() (*Job, bool) {
	pub := atomic.LoadUint64(&rb.publisherCursor)
	sub := atomic.LoadUint64(&rb.subscriberCursor)

	if sub >= pub {
		return nil, false // No data
	}

	// Read job from ring buffer (no synchronization needed - SPSC)
	job := &rb.data[sub&rb.mask]

	// Update subscriber cursor atomically
	atomic.StoreUint64(&rb.subscriberCursor, sub+1)

	return job, true
}

// SubscribeBatch retrieves up to n jobs at once
func (rb *RingBuffer) SubscribeBatch(maxCount int) []Job {
	pub := atomic.LoadUint64(&rb.publisherCursor)
	sub := atomic.LoadUint64(&rb.subscriberCursor)

	available := pub - sub
	if available == 0 {
		return nil
	}

	if uint64(maxCount) > available {
		maxCount = int(available)
	}

	result := make([]Job, maxCount)
	for i := 0; i < maxCount; i++ {
		result[i] = rb.data[(sub+uint64(i))&rb.mask]
	}

	atomic.StoreUint64(&rb.subscriberCursor, sub+uint64(maxCount))

	return result
}

// Reset clears the ring buffer
func (rb *RingBuffer) Reset() {
	atomic.StoreUint64(&rb.publisherCursor, 0)
	atomic.StoreUint64(&rb.subscriberCursor, 0)
}

// Capacity returns the capacity of the ring buffer
func (rb *RingBuffer) Capacity() uint64 {
	return rb.capacity
}
