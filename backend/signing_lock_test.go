package backend

import (
	"github.com/google/uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLockAndUnlock(t *testing.T) {
	logicalStore := &logical.InmemStorage{}

	lock := &DBLock{id:uuid.New(), storage:logicalStore}

	for i := 0 ; i < 10 ; i++ {
		require.NoError(t, lock.Lock())

		isLocked, err := lock.IsLocked()
		require.NoError(t, err)
		require.True(t, isLocked)

		require.NoError(t, lock.UnLock())
	}
}

func TestLockAndRelock(t *testing.T) {
	logicalStore := &logical.InmemStorage{}

	lock := &DBLock{id:uuid.New(), storage:logicalStore}
	require.NoError(t, lock.Lock())

	for i := 0 ; i < 10 ; i++ {
		require.NotNil(t, lock.Lock())
		require.EqualError(t, lock.Lock(), "locked")

		isLocked, err := lock.IsLocked()
		require.NoError(t, err)
		require.True(t, isLocked)
	}
}
