#!/bin/bash

mkdir /data

VAULT_SERVER_CONFIG=/vault/config/vault-config.json
VAULT_SERVER_SCHEMA=http

# Generate SSL certificates if external IP address is provided
if [ "$VAULT_EXTERNAL_ADDRESS" != "" ]; then
  export VAULT_CACERT=/vault/config/ca.pem
  VAULT_SERVER_CONFIG=/vault/config/vault-config-tls.json
  VAULT_SERVER_SCHEMA=https

  /bin/sh /vault/config/vault-tls.sh
fi

export VAULT_SERVER_SCHEMA=$VAULT_SERVER_SCHEMA
export VAULT_ADDR=$VAULT_SERVER_SCHEMA://127.0.0.1:8200
export VAULT_API_ADDR=$VAULT_SERVER_SCHEMA://127.0.0.1:8200

# Start vault server
vault server -config=$VAULT_SERVER_CONFIG -log-level=debug > /data/logs 2&1

sleep 5
if [ "$UNSEAL" = "true" ]; then
  /bin/sh /vault/config/vault-init.sh
  /bin/sh /vault/config/vault-unseal.sh
  /bin/sh /vault/config/vault-plugin.sh
fi

sleep 356000d
