#!/bin/sh

export SHASUM256=$(sha256sum "/vault/plugins/ethsign" | cut -d' ' -f1)

# Register plugin
vault plugin register \
    -sha256=${SHASUM256} \
    -args=--log-format=${LOG_FORMAT} \
    -args=--log-dsn=${LOG_DSN} \
    -args=--log-levels=${LOG_LEVELS} \
    secret ethsign

# Upgrade Ethereum 2.0 Signing Plugin
#vault write /sys/plugins/catalog/secret/ethsign \
#    sha_256=${SHASUM256} \
#    args="--log-format=${LOG_FORMAT}" \
#    args="--log-dsn=${LOG_DSN}" \
#    args="--log-levels=${LOG_LEVELS}" \
#    command=ethsign > /dev/null 2>&1

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

# Enable launchtest network
echo "Enabling Launch Test network"
vault secrets enable \
    -path=ethereum/launchtest \
    -description="Eth Signing Wallet - Launch Test Network" \
    -plugin-name=ethsign plugin > /dev/null 2>&1

echo "Configuring Launch Test network..."
vault write ethereum/launchtest/config \
    network="launchtest" \
    genesis_time="$LAUNCHTESTNET_GENESIS_TIME"

# Reload plugin
curl --header "X-Vault-Token: $(cat /data/keys/vault.root.token)" --request PUT --data '{"plugin": "ethsign"}'  http://127.0.0.1:8200/v1/sys/plugins/reload/backend
