package tests

import (
	"testing"
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
