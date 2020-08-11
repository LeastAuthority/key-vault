package tests

import (
	"fmt"
	"testing"

	"github.com/bloxapp/KeyVault/core"
	"github.com/stretchr/testify/require"
)

type E2E interface {
	Name() string
	Run(t *testing.T)
}

var tests = []E2E{
	&AttestationSigning{},
	&AttestationSigningAccountNotFound{},
	&AttestationDoubleSigning{},
	&AttestationConcurrentSigning{},
}

func TestE2E(t *testing.T) {
	for _, tst := range tests {
		t.Run(tst.Name(), func(t *testing.T) {
			tst.Run(t)
		})
	}
}

func TestNewSeed(t *testing.T) {
	entropy, err := core.GenerateNewEntropy()
	require.NoError(t, err)

	seed, err := core.SeedFromEntropy(entropy, "test_password")
	require.NoError(t, err)

	fmt.Println(string(seed))
}
