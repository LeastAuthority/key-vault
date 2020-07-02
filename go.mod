module github.com/bloxapp/vault-plugin-secrets-eth2.0

go 1.14

require (
	github.com/bloxapp/KeyVault v0.0.0-20200702084557-9df8b01ada4c
	github.com/bloxapp/prysm v0.3.10
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/vault/api v1.0.4
	github.com/hashicorp/vault/sdk v0.1.13
	github.com/herumi/bls-eth-go-binary v0.0.0-20200624084043-9b7da5962ccb // indirect
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/prysm v0.3.10
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/wealdtech/eth2-signer-api v1.3.0
	github.com/wealdtech/go-eth2-types/v2 v2.4.2
	gopkg.in/urfave/cli.v2 v2.0.0-00010101000000-000000000000 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace gopkg.in/urfave/cli.v2 => github.com/urfave/cli/v2 v2.1.1

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20200530091827-df74fa9e9621

replace github.com/herumi/bls-eth-go-binary => github.com/herumi/bls-eth-go-binary v0.0.0-20200605082007-3a76b4c6c599

replace github.com/prysmaticlabs/prysm => github.com/prysmaticlabs/prysm v1.0.0-alpha.10
