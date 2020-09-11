package store

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

// To make sure Logical implements logical.Storage interface.
var _ logical.Storage = &logicalStore{}

// logicalStore wraps logical.Storage to separate storages based on the subnet
type logicalStore struct {
	base   logical.Storage
	prefix string
}

// NewLogical is the constructor of logicalStore.
func NewLogical(base logical.Storage, prefix string) logical.Storage {
	return &logicalStore{
		base:   base,
		prefix: prefix,
	}
}

// List implements logical.Storage interface.
func (s *logicalStore) List(ctx context.Context, key string) ([]string, error) {
	return s.base.List(ctx, s.buildKey(key))
}

// Get implements logical.Storage interface.
func (s *logicalStore) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	return s.base.Get(ctx, s.buildKey(key))
}

// Put implements logical.Storage interface.
func (s *logicalStore) Put(ctx context.Context, entry *logical.StorageEntry) error {
	entry.Key = s.buildKey(entry.Key)
	return s.base.Put(ctx, entry)
}

// Delete implements logical.Storage interface.
func (s *logicalStore) Delete(ctx context.Context, key string) error {
	return s.base.Delete(ctx, s.buildKey(key))
}

func (s *logicalStore) buildKey(key string) string {
	return fmt.Sprintf("%s/%s", s.prefix, key)
}
