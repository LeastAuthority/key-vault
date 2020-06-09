#!/bin/sh
# Init Vault
echo "Init vault..."
vault operator init  &> /vault/keys.txt

if ! grep -q "Vault is already initialized" /vault/keys.txt; then

# Unseal Vault
echo "Unseal vault..."
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 1:' /vault/keys.txt | awk '{print $NF}')
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 2:' /vault/keys.txt | awk '{print $NF}')
vault operator unseal -address=${VAULT_ADDR} $(grep 'Key 3:' /vault/keys.txt | awk '{print $NF}')
vault login $(grep 'Initial Root Token:' /vault/keys.txt | awk '{print $NF}') > /vault/token.txt
echo "###Please save the token information.###"
cat /vault/token.txt


else
    echo "Vault is already initialized!"
fi
