package kvstore

import "errors"

var (
	ErrKeyNotFound    = errors.New("key not found")
	ErrNodeNotFound   = errors.New("node not found")
	ErrNodeDead       = errors.New("node is dead")
	ErrNoAliveNodes   = errors.New("no alive nodes available")
	ErrReplicaFailure = errors.New("insufficient replicas acknowledged write")
	ErrEmptyKey       = errors.New("key must not be empty")
)
