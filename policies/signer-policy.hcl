# Ability to list existing wallet accounts ("list")
path "ethereum/test/accounts" {
  capabilities = ["list"]
}
path "ethereum/launchtest/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/test/accounts/sign-*" {
  capabilities = ["create"]
}
path "ethereum/launchtest/accounts/sign-*" {
  capabilities = ["create"]
}
