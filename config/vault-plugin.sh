#!/bin/sh

export SHASUM256=$(sha256sum "/vault/plugins/ethsign" | cut -d' ' -f1)

# Register plugin
vault plugin register \
    -sha256=${SHASUM256} \
    -args=--log-format=${LOG_FORMAT} \
    -args=--log-dsn=${LOG_DSN} \
    -args=--log-levels=${LOG_LEVELS} \
    secret ethsign

# Enable test network
echo "Enabling Test network..."
vault secrets enable \
    -path=ethereum/test \
    -description="Eth Signing Wallet - Test Network" \
    -plugin-name=ethsign plugin > /dev/null 2>&1

echo "Configuring Test network..."
vault write ethereum/test/config \
    network="test" \
    genesis_time="$TESTNET_GENESIS_TIME"

# Enable zinken network
echo "Enabling Zinken Test network"
vault secrets enable \
    -path=ethereum/zinken \
    -description="Eth Signing Wallet - Zinken Test Network" \
    -plugin-name=ethsign plugin > /dev/null 2>&1

echo "Configuring Zinken Test network..."
vault write ethereum/zinken/config \
    network="zinken" \
    genesis_time="$ZINKEN_GENESIS_TIME"

# Reload plugin
curl --insecure --header "X-Vault-Token: $(cat /data/keys/vault.root.token)" \
        --request PUT \
        --data '{"plugin": "ethsign"}' \
         ${VAULT_SERVER_SCHEMA:-http}://127.0.0.1:8200/v1/sys/plugins/reload/backend
