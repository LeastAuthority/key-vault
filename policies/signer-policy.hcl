# Ability to list existing wallet accounts ("list")
path "ethereum/wallet/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/wallet/accounts/+/sign-*" {
  capabilities = ["create"]
}
