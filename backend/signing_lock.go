package backend

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

// DBLock implements DB slocking mechanism.
type DBLock struct {
	id      uuid.UUID
	storage logical.Storage
}

// NewDBLock is the constructor of DBLock.
func NewDBLock(id uuid.UUID, storage logical.Storage) *DBLock {
	return &DBLock{
		id:      id,
		storage: storage,
	}
}

// Lock locks the DB.
func (lock *DBLock) Lock() error {
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
		Key:      lock.key(),
		Value:    []byte("1"),
		SealWrap: false,
	}
	return lock.storage.Put(context.Background(), entry)
}

// UnLock unlocks the DB.
func (lock *DBLock) UnLock() error {
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

// IsLocked returns true if the DB is locked.
func (lock *DBLock) IsLocked() (bool, error) {
	entry, err := lock.storage.Get(context.Background(), lock.key())
	if err != nil {
		return true, err
	}

	return entry != nil, err
}

func (lock *DBLock) key() string {
	return fmt.Sprintf("lock/%s", lock.id.String())
}
