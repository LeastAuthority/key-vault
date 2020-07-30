# vault-plugin-secrets-eth2.0

## How to use this project?

  1. Build the images and run the containers:

        ```sh
        $ docker-compose up --build
        ```

  2. Execute the container

        ```sh
        $ docker-compose exec vault bash
        ```

  3. Initialize the server
 
        ```sh
        $ vault operator init
        ```

  4. Now we need to unseal the server. A sealed server can't accept any requests. Since vault was initialized with 5 key
  shares and requires a minimum of 3 keys to reconstruct master key, we need to send 3 unseal requests with these
  3 keys. After 3 unseal keys are provided, Sealed in the server status turns to false, and the server is ready to be
  further configured.
 
        ```sh
        $ vault operator unseal
        ```


## Install Ethereum 2.0 Signing Plugin

  1. Login into the vault using root token

        ```sh
        $ vault login Your_Initial_Root_Token
        ```
      
  2. Calculate checksum of the binary from the build

        ```sh
        $ export SHASUM256=$(sha256sum "/vault/plugins/ethsign" | cut -d' ' -f1)
        ```
  
  3. Register the binary with vault server.

        ```sh
        $ vault write /sys/plugins/catalog/secret/ethsign sha_256=${SHASUM256} command=ethsign
        ```

  4. Enable the plugin
  
        ```sh
        $ vault secrets enable -path=ethereum -description="Eth Signing Wallet" -plugin-name=ethsign plugin
        ```
     
## Upgrade Ethereum 2.0 Signing Plugin

  1. Login into the vault using root token

        ```sh
        $ vault login Your_Initial_Root_Token
        ```
      
  2. Calculate checksum of the binary from the build

        ```sh
        $ export SHASUM256=$(sha256sum "/vault/plugins/ethsign" | cut -d' ' -f1)
        ```
  
  3. Register the binary with vault server.

        ```sh
        $ vault write /sys/plugins/catalog/secret/ethsign sha_256=${SHASUM256} command=ethsign
        ```

  4. Trigger a plugin reload
  
        ```sh
        Sample Payload:
        {
            "plugin": "ethsign"
        }
     
        $ curl \
            --header "X-Vault-Token: ..." \
            --request PUT \
            --data @payload.json \
            http://127.0.0.1:8200/v1/sys/plugins/reload/backend
        ```
     

### LIST ACCOUNTS

This endpoint will list all accounts of key-vault stores at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/accounts`  | `200 application/json` |


#### Sample Response

The example below shows output for a query path of `/ethereum/accounts` when there are 1 account.

```
{
    "request_id": "489790dc-b4bd-54e5-be6e-95a894ffc48c",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "accounts": [
            {
                "id": "9676ef06-d238-49f3-ab50-b3fe9930db0f",
                "name": "account-0",
                "validationPubKey": "8a5df36be5f89f9fe19cabadcbb17babc8c518bcd7fe0095c89f83915ea943343fa7dd3c26d8fb6096bce11fbc1ec7d3",
                "withdrawalPubKey": "887abb059075160ce2556a8bfef745898ee3a11b2b6521b09077d422c164929dea277ac8afcacd5b6d729198238f8f6c"
            }
        ]
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}
```


### SIGN ATTESTATION

This endpoint will sign attestation for specific account at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:account_name/sign-attestation`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to sign. This is specified as part of the URL.
* `domain` (`string: <required>`) - Specifies the domain.
* `slot` (`int: <required>`) - Specifies the slot.
* `committeeIndex` (`int: <required>`) - Specifies the committeeIndex.
* `beaconBlockRoot` (`string: <required>`) - Specifies the beaconBlockRoot.
* `sourceEpoch` (`int: <required>`) - Specifies the sourceEpoch.
* `sourceRoot` (`string: <required>`) - Specifies the sourceRoot.
* `targetEpoch` (`int: <required>`) - Specifies the targetEpoch.
* `targetRoot` (`string: <required>`) - Specifies the targetRoot.

#### Sample Response

The example below shows output for the successful sign of `/ethereum/accounts/account1/sign-attestation`.

```
{
    "request_id": "b767dcca-5b10-4a52-1d9a-0a9b81b378ae",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "signature": "kEEOMxNkouz7EOSULfrG6hXzZbIOvRCVVK+lfBofj3U49/PHm7YHji8ac9Gf9vgEFVEmbPp+lhO3OpAElt3yaBajTKaJBWocgXuv64Ojq44tfxLJo6jrzMU5yoP78dYW"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}
```

### SIGN PROPOSAL

This endpoint will sign attestation for specific account at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:account_name/sign-proposal`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to sign. This is specified as part of the URL.
* `domain` (`string: <required>`) - Specifies the domain.
* `slot` (`int: <required>`) - Specifies the slot.
* `proposerIndex` (`int: <required>`) - Specifies the proposerIndex.
* `parentRoot` (`string: <required>`) - Specifies the parentRoot.
* `stateRoot` (`string: <required>`) - Specifies the stateRoot.
* `bodyRoot` (`string: <required>`) - Specifies the bodyRoot.

#### Sample Response

The example below shows output for the successful sign of `/ethereum/accounts/account1/sign-proposal`.

```
{
    "request_id": "b767dcca-5b10-4a52-1d9a-0a9b81b378ae",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "signature": "kPyCp8ID44ceUB3KSp+7QsxqTlGSP2u6/cytr04qJyxkkIKIO/FW57qwH9E7/c48D1PgHsyb8hgoT8/jOLMD7Y/Jt06Qiw80ZRtoS78CzMFYRut/OQot+FzAJcW7Jk0U"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}
```

### SIGN AGGREGATION

This endpoint will sign attestation for specific account at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/accounts/:account_name/sign-aggregation`  | `200 application/json` |

#### Parameters

* `account_name` (`string: <required>`) - Specifies the name of the account to sign. This is specified as part of the URL.
* `domain` (`string: <required>`) - Specifies the domain.
* `dataToSign` (`string: <required>`) - Specifies the slot.

#### Sample Response

The example below shows output for the successful creation of `/ethereum/accounts/account1/sign-aggregation`.

```
{
    "request_id": "b767dcca-5b10-4a52-1d9a-0a9b81b378ae",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "signature": "kPyCp8ID44ceUB3KSp+7QsxqTlGSP2u6/cytr04qJyxkkIKIO/FW57qwH9E7/c48D1PgHsyb8hgoT8/jOLMD7Y/Jt06Qiw80ZRtoS78CzMFYRut/OQot+FzAJcW7Jk0U"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}
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
path "ethereum/accounts/+/sign-*" {
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
path "ethereum/accounts/+/sign-*" {
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
