package tests

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type E2ETest interface {
	Name() string
	Run() error
}

var tests = []E2ETest {
	&AttestationSigning{},
}

func TestE2E(t *testing.T) {
	for _, tst := range tests {
		t.Run(tst.Name(), func(t *testing.T) {
			require.NoError(t, tst.Run())
		})
	}
}