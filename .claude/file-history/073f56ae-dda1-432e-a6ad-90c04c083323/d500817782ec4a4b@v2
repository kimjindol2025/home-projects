package dispatcher

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"grie-engine/internal/circuitbreaker"
	"grie-engine/internal/config"
	"grie-engine/internal/retry"
	"grie-engine/internal/shm"
)

// Worker represents a single Julia worker goroutine
type Worker struct {
	ID            int
	WorkerState   *WorkerState
	juliaProcess  *exec.Cmd
	shmPath       string
	processedJobs uint64 // atomic
	errors        uint64 // atomic
	lastError     atomic.Value // error
	cpuCore       int // CPU core to pin to (0 = no pinning)

	// Phase 3: Operational Excellence
	timeoutCfg     *config.TimeoutConfig
	circuitBreaker *circuitbreaker.CircuitBreaker
	retryConfig    retry.RetryConfig
}

// NewWorker creates a new worker for a Julia instance
func NewWorker(id int, workerState *WorkerState, cpuCore int) *Worker {
	w := &Worker{
		ID:             id,
		WorkerState:    workerState,
		cpuCore:        cpuCore,
		timeoutCfg:     config.DefaultTimeoutConfig(),
		circuitBreaker: circuitbreaker.NewCircuitBreaker(),
		retryConfig:    retry.DefaultRetryConfig(),
	}
	return w
}

// Start initializes the shared memory and launches the Julia process
func (w *Worker) Start(ctx context.Context) error {
	// Create shared memory for this worker
	shmMgr, err := shm.Create(10 * 1024 * 1024) // 10MB per worker
	if err != nil {
		return fmt.Errorf("failed to create shared memory: %w", err)
	}

	w.WorkerState.SHM = shmMgr
	w.shmPath = shmMgr.Path()

	// TODO: Launch Julia process
	// juliaScript := fmt.Sprintf("julia_worker_%d.jl", w.ID)
	// w.juliaProcess = exec.CommandContext(ctx, "julia", juliaScript, w.shmPath)

	return nil
}

// Run executes the worker main loop (pinned to OS thread if cpuCore >= 0)
func (w *Worker) Run(ctx context.Context, jobQueue chan *Job) error {
	// Pin to OS thread if requested
	if w.cpuCore >= 0 {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// TODO: Actually pin to CPU core using taskset
		// cmd := exec.Command("taskset", "-p", "-c", fmt.Sprintf("%d", w.cpuCore), fmt.Sprintf("%d", os.Getpid()))
		// cmd.Run()
	}

	for {
		select {
		case job := <-jobQueue:
			if job == nil {
				// Channel closed
				return nil
			}

			result, err := w.processJob(ctx, job)
			if err != nil {
				atomic.AddUint64(&w.errors, 1)
				w.lastError.Store(err)
				job.Response <- nil
			} else {
				job.Response <- result
			}

			w.WorkerState.IncrementProcessed()

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// processJob sends job to Julia via shared memory and waits for result
// Includes circuit breaker pattern, retry logic, and configurable timeouts
func (w *Worker) processJob(ctx context.Context, job *Job) ([]byte, error) {
	// Check circuit breaker first
	if err := w.circuitBreaker.Allow(); err != nil {
		return nil, &retry.RetryableError{Err: err, Type: retry.ErrorPermanent}
	}

	// Execute with retry logic
	var result []byte
	err := retry.Do(ctx, w.retryConfig, func() error {
		innerErr := w.doProcessJob(ctx, job.Payload)
		if innerErr != nil {
			w.circuitBreaker.RecordFailure()
			return innerErr
		}

		// Read result after successful write
		var readErr error
		result, readErr = w.readJobResult(ctx)
		if readErr != nil {
			w.circuitBreaker.RecordFailure()
			return readErr
		}

		w.circuitBreaker.RecordSuccess()
		return nil
	})

	return result, err
}

// doProcessJob writes job payload to shared memory
func (w *Worker) doProcessJob(ctx context.Context, payload []byte) error {
	if w.WorkerState.SHM == nil {
		return fmt.Errorf("shared memory not initialized")
	}

	// Create context with timeout for SHM write
	writeCtx, cancel := context.WithTimeout(ctx, w.timeoutCfg.SHMWriteTimeout)
	defer cancel()

	// Write job payload to shared memory
	done := make(chan error, 1)
	go func() {
		done <- w.WorkerState.SHM.WriteData(payload)
	}()

	select {
	case err := <-done:
		if err != nil {
			return &retry.RetryableError{Err: fmt.Errorf("failed to write to shared memory: %w", err), Type: retry.ErrorTransient}
		}
		return nil
	case <-writeCtx.Done():
		return &retry.RetryableError{Err: fmt.Errorf("SHM write timeout"), Type: retry.ErrorTransient}
	}
}

// readJobResult polls shared memory for job result with timeout
func (w *Worker) readJobResult(ctx context.Context) ([]byte, error) {
	deadline := time.Now().Add(w.timeoutCfg.JobProcessTimeout)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		default:
			if time.Now().After(deadline) {
				return nil, &retry.RetryableError{Err: fmt.Errorf("timeout waiting for Julia response"), Type: retry.ErrorTransient}
			}

			// Try to read result
			result, err := w.WorkerState.SHM.ReadData(int(w.timeoutCfg.SHMReadTimeout.Milliseconds()))
			if err == nil {
				return result, nil
			}

			// Not ready yet, keep polling
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// GetStats returns worker statistics
func (w *Worker) GetStats() map[string]interface{} {
	lastErr := w.lastError.Load()
	errStr := ""
	if lastErr != nil {
		errStr = lastErr.(error).Error()
	}

	return map[string]interface{}{
		"id":               w.ID,
		"state":            w.WorkerState.state.Load(),
		"processed_jobs":   atomic.LoadUint64(&w.processedJobs),
		"errors":           atomic.LoadUint64(&w.errors),
		"last_error":       errStr,
		"last_seen":        w.WorkerState.GetLastSeen(),
		"shm_path":         w.shmPath,
		"cpu_core":         w.cpuCore,
	}
}

// Close cleans up the worker
func (w *Worker) Close() error {
	if w.WorkerState.SHM != nil {
		if err := w.WorkerState.SHM.Close(); err != nil {
			return fmt.Errorf("failed to close shared memory: %w", err)
		}
	}

	if w.juliaProcess != nil {
		if err := w.juliaProcess.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill Julia process: %w", err)
		}
	}

	return nil
}

// WorkerPool manages multiple workers
type WorkerPool struct {
	workers      []*Worker
	jobQueues    []chan *Job
	dispatcher   *Dispatcher
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	totalStarted uint64 // atomic
	totalStopped uint64 // atomic
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(dispatcher *Dispatcher, workerCount int) (*WorkerPool, error) {
	workers := make([]*Worker, workerCount)
	jobQueues := make([]chan *Job, workerCount)

	for i := 0; i < workerCount; i++ {
		jobQueues[i] = make(chan *Job, 100) // Buffered job queue per worker
		workers[i] = NewWorker(i, dispatcher.workers[i], i%8) // Pin to 8 cores
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		workers:    workers,
		jobQueues:  jobQueues,
		dispatcher: dispatcher,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// Start launches all workers
func (wp *WorkerPool) Start(ctx context.Context) error {
	for i, worker := range wp.workers {
		if err := worker.Start(ctx); err != nil {
			return fmt.Errorf("failed to start worker %d: %w", i, err)
		}

		wp.wg.Add(1)
		go func(w *Worker, idx int) {
			defer wp.wg.Done()
			atomic.AddUint64(&wp.totalStarted, 1)

			if err := w.Run(ctx, wp.jobQueues[idx]); err != nil {
				w.lastError.Store(err)
			}

			atomic.AddUint64(&wp.totalStopped, 1)
		}(worker, i)
	}

	return nil
}

// Submit sends a job to a specific worker
func (wp *WorkerPool) Submit(workerIdx int, job *Job) error {
	if workerIdx < 0 || workerIdx >= len(wp.jobQueues) {
		return fmt.Errorf("invalid worker index: %d", workerIdx)
	}

	select {
	case wp.jobQueues[workerIdx] <- job:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool shutting down")
	default:
		return fmt.Errorf("worker queue full")
	}
}

// Stop gracefully shuts down all workers
func (wp *WorkerPool) Stop() error {
	wp.cancel()

	// Close all job queues
	for _, q := range wp.jobQueues {
		close(q)
	}

	// Wait for all workers to finish
	wp.wg.Wait()

	// Close all workers
	for _, w := range wp.workers {
		if err := w.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Stats returns stats for all workers
func (wp *WorkerPool) Stats() []map[string]interface{} {
	stats := make([]map[string]interface{}, len(wp.workers))
	for i, w := range wp.workers {
		stats[i] = w.GetStats()
	}
	return stats
}
