Stoyan Adasha, [18.06.20 18:55]
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

sleep 356000d
