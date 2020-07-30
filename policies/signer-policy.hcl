# Ability to list existing wallet accounts ("list")
path "ethereum/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/accounts/+/sign-*" {
  capabilities = ["create"]
}
