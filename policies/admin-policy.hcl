# Ability to list existing wallets ("list")
path "ethereum/wallets" {
  capabilities = ["list"]
}

# Ability to create wallets ("create")
path "ethereum/wallets/+" {
  capabilities = ["create"]
}

# Ability to list existing accounts ("list")
path "ethereum/wallets/+/accounts/" {
  capabilities = ["list"]
}

# Ability to create create accounts ("create") and read existing account ("read")
path "ethereum/wallets/+/accounts/+" {
  capabilities = ["create", "read"]
}

# Ability to read deposit data ("read")
path "ethereum/wallets/+/+/+/deposit-data/" {
  capabilities = ["read"]
}

# Ability to sign data ("create")
path "ethereum/wallets/+/+/+/sign-*" {
  capabilities = ["create"]
}

# Ability to export seed ("read")
path "ethereum/export" {
  capabilities = ["read"]
}

# Ability to import seed ("create")
path "ethereum/import" {
  capabilities = ["create"]
}