package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

// ServerConfig holds monitoring server configuration
type ServerConfig struct {
	Addr            string        // Address to listen on (default: ":9090")
	ReadTimeout     time.Duration // Read timeout (default: 5s)
	WriteTimeout    time.Duration // Write timeout (default: 10s)
	ShutdownTimeout time.Duration // Graceful shutdown timeout (default: 5s)
	Registry        *Registry     // Metrics registry (default: DefaultRegistry)
}

// DefaultServerConfig returns standard configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Addr:            ":9090",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 5 * time.Second,
		Registry:        DefaultRegistry,
	}
}

// Server implements the HTTP monitoring server
type Server struct {
	cfg    ServerConfig
	srv    *http.Server
	ln     net.Listener
	closed bool
}

// NewServer creates a new monitoring server
func NewServer(cfg ServerConfig) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":9090"
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 5 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 5 * time.Second
	}
	if cfg.Registry == nil {
		cfg.Registry = DefaultRegistry
	}

	return &Server{cfg: cfg}
}

// Start starts the monitoring server in the background
func (s *Server) Start() error {
	// Create HTTP router
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		s.handleMetrics(w, r)
	})

	// Health check endpoint (liveness)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		s.handleHealth(w, r)
	})

	// Readiness endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		s.handleReady(w, r)
	})

	// Debug pprof endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	// Create HTTP server
	s.srv = &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      mux,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	// Listen
	ln, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.Addr, err)
	}
	s.ln = ln

	// Start server in background
	go func() {
		if err := s.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("monitoring server error: %v", err)
		}
	}()

	log.Printf("monitoring server started on %s", s.cfg.Addr)
	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	if s.srv == nil || s.closed {
		return nil
	}

	s.closed = true

	// Create context with shutdown timeout if not provided
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
		defer cancel()
	}

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	return nil
}

// handleMetrics handles /metrics endpoint (Prometheus format)
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	_ = s.cfg.Registry.WritePrometheusFormat(w)
}

// handleHealth handles /health endpoint (liveness)
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status": "alive",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(response)
}

// handleReady handles /ready endpoint (readiness)
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// For now, always ready
	// In production, check if dispatcher is running, workers are available, etc.
	ready := true

	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		response := map[string]interface{}{
			"status": "not_ready",
			"time":   time.Now().UTC().Format(time.RFC3339),
		}
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status": "ready",
		"time":   time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(response)
}
