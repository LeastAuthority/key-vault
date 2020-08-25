module github.com/bloxapp/vault-plugin-secrets-eth2.0

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Sirupsen/logrus v1.6.0 // indirect
	github.com/bloxapp/KeyVault v0.2.5
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/frankban/quicktest v1.7.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/vault/api v1.0.4
	github.com/hashicorp/vault/sdk v0.1.13
	github.com/kr/pretty v0.2.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/pierrec/lz4 v2.4.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/wealdtech/eth2-signer-api v1.3.0
	github.com/wealdtech/go-eth2-types/v2 v2.5.0
	github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4 v1.1.0
	github.com/wealdtech/go-eth2-wallet-types/v2 v2.6.0
	golang.org/x/net v0.0.0-20200528225125-3c3fba18258b // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
)

replace gopkg.in/urfave/cli.v2 => github.com/urfave/cli/v2 v2.1.1

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20200530091827-df74fa9e9621

replace github.com/herumi/bls-eth-go-binary => github.com/herumi/bls-eth-go-binary v0.0.0-20200605082007-3a76b4c6c599

replace github.com/prysmaticlabs/prysm => github.com/prysmaticlabs/prysm v1.0.0-alpha.23

replace github.com/Sirupsen/logrus v1.6.0 => github.com/sirupsen/logrus v1.4.1
