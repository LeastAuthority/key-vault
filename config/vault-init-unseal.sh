#!/bin/sh

# Init Vault
echo "Init vault..."
vault operator init  &> /data/keys.txt

if ! grep -q "Vault is already initialized" /vault/data.txt; then

# Unseal Vault
echo "Unseal vault..."
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 1:' /data/keys.txt | awk '{print $NF}')
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 2:' /data/keys.txt | awk '{print $NF}')
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 3:' /data/keys.txt | awk '{print $NF}')
vault login $(grep 'Initial Root Token:' /vault/keys.txt | awk '{print $NF}') > /data/token.txt
echo "###Please save the token information.###"
cat /data/token.txt


else
    echo "Vault is already initialized!"
fi
