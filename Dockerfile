#
# STEP 1: Prepare environment
#
FROM golang:1.14-stretch AS preparer

RUN apt-get update                                                        && \
  DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
    curl git zip unzip wget g++ python gcc jq                                \
  && rm -rf /var/lib/apt/lists/*

RUN go version
RUN python --version

WORKDIR /go/src/github.com/bloxapp/vault-plugin-secrets-eth2.0/
COPY go.mod .
COPY go.sum .
RUN go mod download

#
# STEP 2: Build executable binary
#
FROM preparer AS builder

# Copy files and install app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static -lm"' -o ethsign .

#
# STEP 3: Get vault image and copy the plugin
#
FROM vault:latest AS runner

# Download dependencies
RUN apk -v --update --no-cache add \
    bash ca-certificates

WORKDIR /vault/plugins/

COPY --from=builder /go/src/github.com/bloxapp/vault-plugin-secrets-eth2.0/ethsign ./ethsign
COPY ./config/vault-config.json /vault/config/vault-config.json
COPY ./config/vault-init-unseal.sh /vault/config/vault-init-unseal.sh
COPY ./config/entrypoint.sh /vault/config/entrypoint.sh
RUN chown vault /vault/config/vault-init-unseal.sh 
RUN chown vault /vault/config/entrypoint.sh
RUN apk add jq

WORKDIR /

# Expose port 8200
EXPOSE 8200

# Run vault
CMD ["/vault/config/entrypoint.sh"]
