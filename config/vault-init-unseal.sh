#!/bin/bash
set -e

echo "Initialize Vault"
vault status | tee /data/vault.status > /dev/null


SEALED=$(grep 'Sealed' /data/vault.status | awk '{print $NF}')

if [ "$SEALED" = "false" ]; then

  echo "Vault is already initialized."

else

echo "Initialize Vault"
vault status | tee /data/vault.status > /dev/null
vault operator init | tee /data/vault.init > /dev/null 

cat /data/vault.init

echo "Unsealing Vault"
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 1:' /data/vault.init | awk '{print $NF}') 
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 2:' /data/vault.init | awk '{print $NF}')
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 3:' /data/vault.init | awk '{print $NF}')

echo "Login Vault"
vault login $(grep 'Initial Root Token:' /data/vault.init | awk '{print $NF}') > /data/token.txt

echo "Vault setup complete."

instructions() {
  cat <<EOF

The unseal keys and root token have been stored in /data directory.

  /data/vault.init
  /data/token.txt
EOF

  exit 1
}

instructions

fi
