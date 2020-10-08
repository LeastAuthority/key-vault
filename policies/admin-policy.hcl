# Ability to list existing wallet accounts ("list")
path "ethereum/test/accounts" {
  capabilities = ["list"]
}
path "ethereum/zinken/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/test/accounts/sign-*" {
  capabilities = ["create"]
}
path "ethereum/zinken/accounts/sign-*" {
  capabilities = ["create"]
}

# Ability to update storage ("create")
path "ethereum/test/storage" {
  capabilities = ["create"]
}
path "ethereum/zinken/storage" {
  capabilities = ["create"]
}
