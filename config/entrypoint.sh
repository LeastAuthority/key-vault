#!/bin/sh

vault server -config=/vault/config/vault-config.json > /dev/null 2&1
sleep 5
vault operator init -status
if [[ $? -eq 2 ]]; then
 vault operator init -key-shares=1 -key-threshold=1 -recovery-shares=1 -recovery-threshold=1 -format=json > /tmp/vault.init 2>&1
  cat /tmp/vault.init | jq -r '.root_token' > /keys/vault.root.token && chmod 0400 /keys/vault.root.token 
  cat /tmp/vault.init | jq -r '.unseal_keys_b64[]' > /keys/vault.unseal.token && chmod 0400 /keys/vault.unseal.token
  rm -f /tmp/vault.*
fi 

vault status 
if [[ $? -eq 2 ]]; then 
  vault operator unseal $(cat /keys/vault.unseal.token ) >/dev/null 2>&1
  vault login $(cat /keys/vault.root.token ) >/dev/null 2>&1
fi 

# Upgrade Ethereum 2.0 Signing Plugin
#vault login $(cat /keys/vault.root.token)
#export SHASUM256=$(sha256sum "/vault/plugins/ethsign" | cut -d' ' -f1)
#vault write /sys/plugins/catalog/secret/ethsign sha_256=${SHASUM256} command=ethsign
#curl --header "X-Vault-Token: $(cat /keys/vault.root.token)" --request PUT --data '{"plugin": "ethsign"}'  http://127.0.0.1:8200/v1/sys/plugins/reload/backend

sleep 356000d
