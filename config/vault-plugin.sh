#!/bin/sh

# Upgrade Ethereum 2.0 Signing Plugin
export SHASUM256=$(sha256sum "/vault/plugins/ethsign" | cut -d' ' -f1)
vault write /sys/plugins/catalog/secret/ethsign sha_256=${SHASUM256} command=ethsign > /dev/null 2>&1

# Enable test network
vault secrets enable \
    -path=ethereum/test \
    -description="Eth Signing Wallet - Test Network" \
    -plugin-name=ethsign plugin > /dev/null 2>&1

# Enable launchtest network
vault secrets enable \
    -path=ethereum/launchtest \
    -description="Eth Signing Wallet - Launch Test Network" \
    -plugin-name=ethsign plugin > /dev/null 2>&1

# Reload plugin
curl --header "X-Vault-Token: $(cat /data/keys/vault.root.token)" --request PUT --data '{"plugin": "ethsign"}'  http://127.0.0.1:8200/v1/sys/plugins/reload/backend
