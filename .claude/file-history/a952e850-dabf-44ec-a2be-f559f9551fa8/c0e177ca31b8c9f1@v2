package monitoring

import (
	"bytes"
	"testing"
)

func TestCounterInc(t *testing.T) {
	c := NewCounter("test", "test counter")
	if c.Value() != 0 {
		t.Errorf("expected initial value 0, got %d", c.Value())
	}

	c.Inc()
	if c.Value() != 1 {
		t.Errorf("expected value 1, got %d", c.Value())
	}

	c.Add(5)
	if c.Value() != 6 {
		t.Errorf("expected value 6, got %d", c.Value())
	}
}

func TestGaugeSet(t *testing.T) {
	g := NewGauge("test", "test gauge")
	if g.Value() != 0 {
		t.Errorf("expected initial value 0, got %d", g.Value())
	}

	g.Set(10)
	if g.Value() != 10 {
		t.Errorf("expected value 10, got %d", g.Value())
	}

	g.Inc()
	if g.Value() != 11 {
		t.Errorf("expected value 11, got %d", g.Value())
	}

	g.Dec()
	if g.Value() != 10 {
		t.Errorf("expected value 10, got %d", g.Value())
	}

	g.Add(-5)
	if g.Value() != 5 {
		t.Errorf("expected value 5, got %d", g.Value())
	}
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	// Test job metrics
	r.JobsReceived.Inc()
	r.JobsDropped.Add(2)
	r.JobsProcessed.Inc()
	r.JobsFailed.Inc()

	if r.JobsReceived.Value() != 1 {
		t.Errorf("expected JobsReceived=1, got %d", r.JobsReceived.Value())
	}
	if r.JobsDropped.Value() != 2 {
		t.Errorf("expected JobsDropped=2, got %d", r.JobsDropped.Value())
	}
	if r.JobsProcessed.Value() != 1 {
		t.Errorf("expected JobsProcessed=1, got %d", r.JobsProcessed.Value())
	}

	// Test gauge metrics
	r.QueueSize.Set(100)
	r.ActiveWorkers.Set(5)

	if r.QueueSize.Value() != 100 {
		t.Errorf("expected QueueSize=100, got %d", r.QueueSize.Value())
	}
	if r.ActiveWorkers.Value() != 5 {
		t.Errorf("expected ActiveWorkers=5, got %d", r.ActiveWorkers.Value())
	}
}

func TestPrometheusFormat(t *testing.T) {
	r := NewRegistry()
	r.JobsReceived.Add(10)
	r.JobsProcessed.Add(5)

	var buf bytes.Buffer
	if err := r.WritePrometheusFormat(&buf); err != nil {
		t.Fatalf("WritePrometheusFormat failed: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}

	// Check for required headers
	if !bytes.Contains(buf.Bytes(), []byte("# HELP grie_jobs_received_total")) {
		t.Error("missing HELP for grie_jobs_received_total")
	}
	if !bytes.Contains(buf.Bytes(), []byte("# TYPE grie_jobs_received_total counter")) {
		t.Error("missing TYPE for grie_jobs_received_total")
	}

	// Check for values
	if !bytes.Contains(buf.Bytes(), []byte("grie_jobs_received_total 10")) {
		t.Error("missing or incorrect value for grie_jobs_received_total")
	}
	if !bytes.Contains(buf.Bytes(), []byte("grie_jobs_processed_total 5")) {
		t.Error("missing or incorrect value for grie_jobs_processed_total")
	}
}

func TestStatsJSON(t *testing.T) {
	r := NewRegistry()
	r.JobsReceived.Add(10)
	r.JobsProcessed.Add(5)
	r.QueueSize.Set(50)

	stats := r.StatsJSON()

	jobStats := stats["jobs"].(map[string]interface{})
	if jobStats["received"].(uint64) != 10 {
		t.Errorf("expected jobs.received=10, got %v", jobStats["received"])
	}
	if jobStats["processed"].(uint64) != 5 {
		t.Errorf("expected jobs.processed=5, got %v", jobStats["processed"])
	}

	queueStats := stats["queue"].(map[string]interface{})
	if queueStats["size"].(int64) != 50 {
		t.Errorf("expected queue.size=50, got %v", queueStats["size"])
	}
}

func TestConcurrentMetrics(t *testing.T) {
	c := NewCounter("test", "test counter")
	done := make(chan bool, 10)

	// 10 goroutines incrementing the counter
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			for j := 0; j < 100; j++ {
				c.Inc()
			}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 1000 increments
	if c.Value() != 1000 {
		t.Errorf("expected counter=1000, got %d", c.Value())
	}
}
