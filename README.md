# vault-plugin-secrets-eth2.0

## How to use this project?

  1. Build the images and run the containers:

        ```sh
        $ docker-compose up -d --build
        ```

  2. Execute the container

        ```sh
        $ docker-compose exec vault bash
        ```

  3. Read the root token
 
        ```sh
        $ docker-compose exec -T vault cat /data/keys/vault.root.token
        ```

  4. Initialize the server
 
        ```sh
        $ vault operator init
        ```

  5. Now we need to unseal the server. A sealed server can't accept any requests. Since vault was initialized with 5 key
  shares and requires a minimum of 3 keys to reconstruct master key, we need to send 3 unseal requests with these
  3 keys. After 3 unseal keys are provided, Sealed in the server status turns to false, and the server is ready to be
  further configured.
 
        ```sh
        $ vault operator unseal
        ```

## Access Policies
The plugin's endpoint paths are designed such that admin-level access policies vs. signer-level access policies can be easily separated.

### Sample Signer Level Policy:
Use the following policy to assign to a signer level access token, with the abilities to list accounts and sign transactions.

```
# Ability to list existing wallet accounts ("list")
path "ethereum/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/accounts/sign-*" {
  capabilities = ["create"]
}
```

### Sample Admin Level Policy:
Use the following policy to assign to a admin level access token, with the full ability to update storage, list accounts and sign transactions.

```
# Ability to list existing wallet accounts ("list")
path "ethereum/accounts" {
  capabilities = ["list"]
}

# Ability to sign data ("create")
path "ethereum/accounts/sign-*" {
  capabilities = ["create"]
}

# Ability to update storage ("create")
path "ethereum/storage" {
  capabilities = ["create"]
}
```

## How to use policies?

  1. Create a new policy named admin:

        ```sh
        $ vault policy write admin policies/admin-policy.hcl
        ```

  2. Create a token attached to admin policy:

        ```sh
        $ vault token create -policy="admin"
        ```

  3. Create a new policy named signer:
 
        ```sh
        $ vault policy write signer policies/signer-policy.hcl
        ```

  4. Create a token attached to signer policy:

        ```sh
        $ vault token create -policy="signer"
        ```

## About testing

There are 2 types of tests in the project: end-to-end and unit ones. 
In order to run all tests including e2e ones you will need to do the following command:
```bash
$ make test
``` 

New e2e tests should be placed in `./e2e/tests` directory and implement `E2E` interface.
Use the current format to add new tests. 