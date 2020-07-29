package backend

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

type DBLock struct {
	id uuid.UUID
	storage logical.Storage
}

func (lock *DBLock) Lock () error {
	// if locked return error
	locked, err := lock.IsLocked()
	if err != nil {
		return err
	}
	if locked {
		return fmt.Errorf("locked")
	}

	// add lock to db
	entry := &logical.StorageEntry{
		Key:     lock.key(),
		Value:    []byte("1"),
		SealWrap: false,
	}
	return lock.storage.Put(context.Background(), entry)
}

func (lock *DBLock) UnLock () error {
	// check if locked
	locked, err := lock.IsLocked()
	if err != nil {
		return err
	}
	if !locked {
		return nil
	}

	// if not, unlock
	return lock.storage.Delete(context.Background(), lock.key())
}

func (lock *DBLock) IsLocked() (bool, error) {
	entry, err := lock.storage.Get(context.Background(), lock.key())
	if err != nil {
		return true, err
	}

	return entry != nil, err
}

func (lock *DBLock)key () string {
	return fmt.Sprintf("lock/%s", lock.id.String())
}