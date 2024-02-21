package adapters

import (
	"context"
	"sync"
	"time"
)

type RateLimitMemoryStorageAdapter struct {
	mutexAccesses sync.Mutex
	mutexBlocks   sync.Mutex
	accesses      map[string]*map[string]*[]*time.Time
	blocks        map[string]*map[string]*time.Time
}

func NewRateLimitMemoryStorageAdapter() *RateLimitMemoryStorageAdapter {
	adapter := RateLimitMemoryStorageAdapter{}
	adapter.mutexAccesses = sync.Mutex{}
	adapter.mutexBlocks = sync.Mutex{}
	adapter.accesses = map[string]*map[string]*[]*time.Time{}
	adapter.blocks = map[string]*map[string]*time.Time{}
	return &adapter
}

func (s *RateLimitMemoryStorageAdapter) IncrementAccesses(ctx context.Context, keyType string, key string, maxAccesses int64) (bool, int64, error) {
	s.mutexAccesses.Lock()
	defer s.mutexAccesses.Unlock()

	keyTypeData, ok := s.accesses[keyType]
	if !ok {
		keyTypeData = &map[string]*[]*time.Time{}
		s.accesses[keyType] = keyTypeData
	}

	keyData, ok := (*keyTypeData)[key]
	if !ok {
		keyData = &[]*time.Time{}
		(*keyTypeData)[key] = keyData
	}

	filteredKeyData, count := s.filterInLastSecond(keyData)

	if count >= maxAccesses {
		return false, count, nil
	}

	now := time.Now()
	updatedKeyData := append(*filteredKeyData, &now)
	(*keyTypeData)[key] = &updatedKeyData

	return true, count + 1, nil
}

func (s *RateLimitMemoryStorageAdapter) filterInLastSecond(keyData *[]*time.Time) (*[]*time.Time, int64) {
	now := time.Now()
	filtered := []*time.Time{}

	for _, value := range *keyData {
		if now.Sub(*value).Seconds() < 1 {
			filtered = append(filtered, value)
		}
	}

	return &filtered, int64(len(filtered))
}

func (s *RateLimitMemoryStorageAdapter) GetBlock(ctx context.Context, keyType string, key string) (*time.Time, error) {
	s.mutexBlocks.Lock()
	defer s.mutexBlocks.Unlock()

	keyTypeData, ok := s.blocks[keyType]
	if !ok {
		return nil, nil
	}

	blockedUntil, ok := (*keyTypeData)[key]
	if !ok {
		return nil, nil
	}

	if blockedUntil.After(time.Now()) {
		return blockedUntil, nil
	}

	delete(*keyTypeData, key)
	return nil, nil
}

func (s *RateLimitMemoryStorageAdapter) AddBlock(ctx context.Context, keyType string, key string, milliseconds int64) (*time.Time, error) {
	s.mutexBlocks.Lock()
	defer s.mutexBlocks.Unlock()

	keyTypeData, ok := s.blocks[keyType]
	if !ok {
		keyTypeData = &map[string]*time.Time{}
		s.blocks[keyType] = keyTypeData
	}

	blockedUntil := time.Now().Add(time.Duration(int64(time.Millisecond) * milliseconds))
	(*keyTypeData)[key] = &blockedUntil

	return &blockedUntil, nil
}
