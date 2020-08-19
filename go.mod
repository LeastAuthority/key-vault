module github.com/bloxapp/vault-plugin-secrets-eth2.0

go 1.14

require (
	github.com/bloxapp/KeyVault v0.2.4
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/vault/api v1.0.4
	github.com/hashicorp/vault/sdk v0.1.13
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/prysm v0.3.10
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.1
	github.com/wealdtech/eth2-signer-api v1.3.0
	github.com/wealdtech/go-eth2-types/v2 v2.5.0
	github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4 v1.1.0
	github.com/wealdtech/go-eth2-wallet-types/v2 v2.6.0
)

replace gopkg.in/urfave/cli.v2 => github.com/urfave/cli/v2 v2.1.1

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20200530091827-df74fa9e9621

replace github.com/herumi/bls-eth-go-binary => github.com/herumi/bls-eth-go-binary v0.0.0-20200605082007-3a76b4c6c599

replace github.com/prysmaticlabs/prysm => github.com/prysmaticlabs/prysm v1.0.0-alpha.10
