package tests

import (
	"testing"
)

func TestAttestationSigning(t *testing.T) {
	test := &AttestationSigning{}
	test.Run(t)
}

func TestAttestationSigningAccountNotFound(t *testing.T) {
	test := &AttestationSigningAccountNotFound{}
	test.Run(t)
}

func TestAttestationDoubleSigning(t *testing.T) {
	test := &AttestationDoubleSigning{}
	test.Run(t)
}

func TestE2E(t *testing.T) {
	test := &AttestationConcurrentSigning{}
	test.Run(t)
}
