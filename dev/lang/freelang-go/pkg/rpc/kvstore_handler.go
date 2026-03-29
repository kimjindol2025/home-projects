package rpc

import (
	"fmt"

	"github.com/kimjindol2025/freelang-go/pkg/kvstore"
)

// KVStoreHandler wraps a *kvstore.Cluster and exposes its operations as RPC methods.
type KVStoreHandler struct {
	cluster *kvstore.Cluster
	codec   Codec
}

// NewKVStoreHandler creates a handler wrapping cluster.
func NewKVStoreHandler(cluster *kvstore.Cluster, codec Codec) *KVStoreHandler {
	return &KVStoreHandler{
		cluster: cluster,
		codec:   codec,
	}
}

// Register installs Set, Get, Delete, Stats on server.
// Method names: "KV.Set", "KV.Get", "KV.Delete", "KV.Stats"
func (h *KVStoreHandler) Register(s *Server) error {
	if err := s.Register("KV.Set", h.handleSet); err != nil {
		return err
	}
	if err := s.Register("KV.Get", h.handleGet); err != nil {
		return err
	}
	if err := s.Register("KV.Delete", h.handleDelete); err != nil {
		return err
	}
	if err := s.Register("KV.Stats", h.handleStats); err != nil {
		return err
	}
	return nil
}

// handleSet handles KV.Set: args are {key, value} (both strings).
func (h *KVStoreHandler) handleSet(args []byte) ([]byte, error) {
	var key string
	var value string

	// Decode key.
	keyBuf := make([]byte, len(args))
	copy(keyBuf, args)
	if err := h.codec.Unmarshal(keyBuf, &key); err != nil {
		return nil, fmt.Errorf("unmarshal key: %w", err)
	}

	// Find where value starts (after key encoding).
	codec := NewBinaryCodec() // to skip key in buffer
	keyEncoded, _ := codec.Marshal(key, make([]byte, 0))
	if len(args) <= len(keyEncoded) {
		return nil, fmt.Errorf("missing value in args")
	}

	valueBuf := args[len(keyEncoded):]
	if err := h.codec.Unmarshal(valueBuf, &value); err != nil {
		return nil, fmt.Errorf("unmarshal value: %w", err)
	}

	// Call cluster.Set.
	if err := h.cluster.Set(key, value); err != nil {
		return nil, err
	}

	// Return nil result.
	return h.codec.Marshal(nil, make([]byte, 0))
}

// handleGet handles KV.Get: args is a string (key), result is a string (value).
func (h *KVStoreHandler) handleGet(args []byte) ([]byte, error) {
	var key string
	if err := h.codec.Unmarshal(args, &key); err != nil {
		return nil, fmt.Errorf("unmarshal key: %w", err)
	}

	// Call cluster.Get.
	value, err := h.cluster.Get(key)
	if err != nil {
		return nil, err
	}

	// Marshal result.
	return h.codec.Marshal(value, make([]byte, 0))
}

// handleDelete handles KV.Delete: args is a string (key), result is a bool.
func (h *KVStoreHandler) handleDelete(args []byte) ([]byte, error) {
	var key string
	if err := h.codec.Unmarshal(args, &key); err != nil {
		return nil, fmt.Errorf("unmarshal key: %w", err)
	}

	// Call cluster.Delete.
	err := h.cluster.Delete(key)
	success := err == nil

	// Marshal result (true if deleted, false otherwise).
	return h.codec.Marshal(success, make([]byte, 0))
}

// handleStats handles KV.Stats: no args, result is a struct with counts.
func (h *KVStoreHandler) handleStats(args []byte) ([]byte, error) {
	stats := h.cluster.Stats()

	// Marshal as four int64 values: TotalNodes, AliveNodes, DeadNodes, TotalKeys.
	result := make([]byte, 0, 256)
	var err error

	result, err = h.codec.Marshal(int64(stats.TotalNodes), result)
	if err != nil {
		return nil, err
	}

	result, err = h.codec.Marshal(int64(stats.AliveNodes), result)
	if err != nil {
		return nil, err
	}

	result, err = h.codec.Marshal(int64(stats.DeadNodes), result)
	if err != nil {
		return nil, err
	}

	result, err = h.codec.Marshal(int64(stats.TotalKeys), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
