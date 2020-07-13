# Ability to list existing accounts ("list")
path "ethereum/wallets/+/accounts/" {
  capabilities = ["list"]
}

# Ability to read existing account ("read")
path "ethereum/wallets/+/accounts/+" {
  capabilities = ["read"]
}

# Ability to sign data ("create")
path "ethereum/wallets/+/+/+/sign-*" {
  capabilities = ["create"]
}
