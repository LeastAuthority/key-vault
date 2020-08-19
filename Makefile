IMAGE_NAME = "vault-plugin-secrets-eth2.0:$(shell git rev-parse --short HEAD)"

test:
	echo $(IMAGE_NAME)
	docker build -t $(IMAGE_NAME) -f Dockerfile .
	VAULT_PLUGIN_IMAGE="$(IMAGE_NAME)" go test -cover -race -p 1 ./...
