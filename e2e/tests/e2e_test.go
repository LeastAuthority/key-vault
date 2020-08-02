package tests

import (
	"testing"
)

type E2ETest interface {
	Name() string
	Run(t *testing.T)
}

var tests = []E2ETest {
	&AttestationSigning{},
	&AttestationDoubleSigning{},
}

func TestE2E(t *testing.T) {
	for _, tst := range tests {
		t.Run(tst.Name(), func(t *testing.T) {
			tst.Run(t)
		})
	}


}