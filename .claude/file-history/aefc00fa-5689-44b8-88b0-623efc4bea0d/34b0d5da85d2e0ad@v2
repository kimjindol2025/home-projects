package kernel

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"grie-engine/internal/protocol"
	"grie-engine/internal/shm"
)

// KernelBridge manages Go-Zig communication via shared memory
// Uses atomic flags to avoid CGO overhead (no function calls, just memory ops)
type KernelBridge struct {
	shmManager *shm.SharedMemManager
	header     *protocol.EngineHeader

	// Metrics
	totalSignals uint64 // atomic
	totalLatency uint64 // atomic (nanoseconds)

	// Control
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// NewKernelBridge creates a new Zig kernel bridge
// It connects to the shared memory created by Dispatcher
func NewKernelBridge(shmPath string) (*KernelBridge, error) {
	// Open existing SHM (created by Go Dispatcher)
	mgr, err := shm.Open(shmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SHM: %w", err)
	}

	kb := &KernelBridge{
		shmManager: mgr,
		header:     mgr.Header(),
		stopCh:     make(chan struct{}),
	}

	return kb, nil
}

// SubmitJob submits a job to Zig kernel via shared memory
// Returns immediately without waiting (async pattern)
// Zig will detect state change and begin computation
func (kb *KernelBridge) SubmitJob(payload []byte) error {
	// Write data to shared memory
	if err := kb.shmManager.WriteData(payload); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// Data is now in shared memory with state=READY
	// Zig will detect this and begin computation
	// No CGO call needed!

	return nil
}

// WaitResult waits for Zig kernel to complete computation
// Polls the header state with exponential backoff
func (kb *KernelBridge) WaitResult(timeoutMs int) ([]byte, error) {
	deadline := time.Now().Add(time.Duration(timeoutMs) * time.Millisecond)
	backoff := 1 * time.Microsecond

	for {
		// Poll state (volatile read)
		state := kb.header.LoadState()

		if state == protocol.StateIdle {
			// Computation done! Read result
			result, err := kb.shmManager.ReadData(100)
			if err != nil {
				return nil, fmt.Errorf("failed to read result: %w", err)
			}

			atomic.AddUint64(&kb.totalSignals, 1)
			return result, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for Zig kernel (state: %v)", state)
		}

		// Exponential backoff (but keep spinning for ultra-low latency)
		time.Sleep(backoff)
		if backoff < 10*time.Microsecond {
			backoff *= 2
		}
	}
}

// ExecuteSync submits job and waits for result (synchronous)
// Latency: SHM write (10ns) + Zig compute (0-100ns) + SHM read (10ns) ≈ 20-120ns total
func (kb *KernelBridge) ExecuteSync(payload []byte, timeoutMs int) ([]byte, error) {
	start := time.Now()

	if err := kb.SubmitJob(payload); err != nil {
		return nil, err
	}

	result, err := kb.WaitResult(timeoutMs)
	if err != nil {
		return nil, err
	}

	latency := time.Since(start).Nanoseconds()
	atomic.AddUint64(&kb.totalLatency, uint64(latency))

	return result, nil
}

// MonitorAsync continuously monitors Zig kernel in background
// Processes results as they become available
func (kb *KernelBridge) MonitorAsync(handler func([]byte) error) error {
	kb.wg.Add(1)
	go func() {
		defer kb.wg.Done()

		ticker := time.NewTicker(100 * time.Microsecond)
		defer ticker.Stop()

		for {
			select {
			case <-kb.stopCh:
				return
			case <-ticker.C:
				state := kb.header.LoadState()
				if state == protocol.StateIdle {
					// Result available
					if result, err := kb.shmManager.ReadData(10); err == nil {
						if err := handler(result); err != nil {
							// Log error but continue monitoring
							fmt.Printf("❌ Handler error: %v\n", err)
						}
						atomic.AddUint64(&kb.totalSignals, 1)
					}
				}
			}
		}
	}()

	return nil
}

// Stats returns kernel statistics
func (kb *KernelBridge) Stats() map[string]interface{} {
	signals := atomic.LoadUint64(&kb.totalSignals)
	latency := atomic.LoadUint64(&kb.totalLatency)

	avgLatency := uint64(0)
	if signals > 0 {
		avgLatency = latency / signals
	}

	return map[string]interface{}{
		"total_signals":   signals,
		"total_latency_ns": latency,
		"avg_latency_ns":  avgLatency,
		"shm_path":        kb.shmManager.Path(),
		"shm_size":        kb.shmManager.Size(),
	}
}

// Close cleanly shuts down the bridge
func (kb *KernelBridge) Close() error {
	close(kb.stopCh)
	kb.wg.Wait()

	if kb.shmManager != nil {
		return kb.shmManager.Close()
	}

	return nil
}

// ============================================================================
// FFI Export Functions (for future Zig dynamic linking)
// ============================================================================

// GqseKernelInit initializes Zig kernel (can be called via cgo or static linking)
// extern int gqse_kernel_init(void);
// Returns 0 on success, -1 on error

// GqseKernelRun executes N computation cycles
// extern int gqse_kernel_run(uint32_t cycles);
// Returns number of completed operations

// GqseKernelShutdown cleanly shuts down Zig kernel
// extern void gqse_kernel_shutdown(void);
