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
     
### LIST WALLETS

This endpoint will list all wallets stores at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/wallets`  | `200 application/json` |

#### Parameters

* `mount-path` (`string: <required>`) - Specifies the path of the wallets to list. This is specified as part of the URL.

#### Sample Response

The example below shows output for a query path of `/ethereum/wallets/` when there are 4 accounts.

```
{
    "request_id": "9c66d2fb-dfae-9902-19d2-c4b91ac1ab72",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "wallets": [
            {
                "id": "d0cf6367-865a-4c80-bfc2-c036ecfa6eda",
                "indexMapper": {
                    "account2": "2a75d4eb-e541-4536-a237-8705a8f36a19"
                },
                "key": {
                    "id": "22cc1c64-2853-4d28-8ad7-519f93098454",
                    "path": "m/12381/3600/0",
                    "pubkey": "9038df960844983756bf40d30e02e5fedd62f59745bae4a3393557d2fcf49be72a22ed3a6b6b325744fecf55aabf986c"
                },
                "name": "wallet1",
                "type": "HD"
            }
        ]
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}

```

### CREATE WALLET

This endpoint will create an Ethereum account at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/wallets/:name`  | `200 application/json` |

#### Parameters

* `name` (`string: <required>`) - Specifies the name of the wallet to create. This is specified as part of the URL.

#### Sample Response

The example below shows output for the successful creation of `/ethereum/wallets/wallet2`.

```
{
    "request_id": "7a1988e7-8d40-cab8-c33b-b81f1869503e",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "wallet": {
            "id": "3486507e-5318-4d26-aef2-e52cc241efbf",
            "indexMapper": {},
            "key": {
                "id": "71143ff7-fef5-4296-89e4-3d856924327c",
                "path": "m/12381/3600/0",
                "pubkey": "b29769ca8ee1c9f299977daf91432b279547079ce26cb98b350733de36deb9c54819ffa3f2a89ef85ec41d5c7b0c50d1"
            },
            "name": "wallet1",
            "type": "HD"
        }
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}
```

### LIST ACCOUNTS

This endpoint will list all accounts of specific wallet stores at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `LIST`  | `:mount-path/wallets/:wallet_name/accounts`  | `200 application/json` |

#### Parameters

* `wallet_name` (`string: <required>`) - Specifies the name of the wallet to get accounts of. This is specified as part of the URL.

#### Sample Response

The example below shows output for a query path of `/ethereum/wallets/wallet1/accounts` when there are 2 accounts.

```
{
    "request_id": "5402df19-ffcb-1969-935b-d76c214462a3",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "accounts": [
            {
                "id": "453fe5cc-7d52-4085-b660-b970f35925b6",
                "key": {
                    "id": "a3635e45-e6f5-41e9-9167-10e5898f522e",
                    "path": "m/12381/3600/0/0/0",
                    "pubkey": "845ebb2be5ed29e332b9d7f1825e9512bbb113cdb1cd536c311e2b19e9dd992973c991cd210ca1dc50d972253cc76307"
                },
                "name": "account2",
                "parentWalletId": "79561078-100d-4c4a-9e35-dfe3a3c34168",
                "type": "Validation"
            }
        ]
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}

```

### CREATE ACCOUNT

This endpoint will create an Ethereum 2.0 account of specific wallet at a path.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `POST`  | `:mount-path/wallets/:wallet_name/accounts/:account_name`  | `200 application/json` |

#### Parameters

* `walet_name` (`string: <required>`) - Specifies the name of the wallet to create account in. This is specified as part of the URL.
* `account_name` (`string: <required>`) - Specifies the name of the account to create. This is specified as part of the URL.

#### Sample Response

The example below shows output for the successful creation of `/ethereum/wallets/wallet1/accounts/account1`.

```
{
    "request_id": "b767dcca-5b10-4a52-1d9a-0a9b81b378ae",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "account": {
            "id": "982be00d-9453-4c43-9239-348afaf8595b",
            "key": {
                "id": "e75e8142-fd54-4baa-b155-fcb6e62e2c84",
                "path": "m/12381/3600/0/0/0",
                "pubkey": "9184a9163413073e6432de3409615014fe2303bcf49f5ebf077ab52c74df77d68c7fbc7499cc7cfff47421f7dad675bf"
            },
            "name": "account1",
            "parentWalletId": "3d809558-385f-4a0f-8f55-769bbb0b5586",
            "type": "Validation"
        }
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
| `POST`  | `:mount-path/wallets/:wallet_name/accounts/:account_name/sign-attestation`  | `200 application/json` |

#### Parameters

* `walet_name` (`string: <required>`) - Specifies the name of the wallet of the account to sign. This is specified as part of the URL.
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

The example below shows output for the successful creation of `/ethereum/wallets/wallet1/accounts/account1/sign-attestation`.

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
| `POST`  | `:mount-path/wallets/:wallet_name/accounts/:account_name/sign-proposal`  | `200 application/json` |

#### Parameters

* `walet_name` (`string: <required>`) - Specifies the name of the wallet of the account to sign. This is specified as part of the URL.
* `account_name` (`string: <required>`) - Specifies the name of the account to sign. This is specified as part of the URL.
* `domain` (`string: <required>`) - Specifies the domain.
* `slot` (`int: <required>`) - Specifies the slot.
* `proposerIndex` (`int: <required>`) - Specifies the proposerIndex.
* `parentRoot` (`string: <required>`) - Specifies the parentRoot.
* `stateRoot` (`string: <required>`) - Specifies the stateRoot.
* `bodyRoot` (`string: <required>`) - Specifies the bodyRoot.

#### Sample Response

The example below shows output for the successful creation of `/ethereum/wallets/wallet1/accounts/account1/sign-proposal`.

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

### GET DEPOSIT DATA

This endpoint will return deposit data.

| Method  | Path | Produces |
| ------------- | ------------- | ------------- |
| `GET`  | `:mount-path/wallets/:wallet_name/accounts/:account_name/deposit-data/`  | `200 application/json` |

#### Parameters

* `walet_name` (`string: <required>`) - Specifies the name of the wallet to create account in. This is specified as part of the URL.
* `account_name` (`string: <required>`) - Specifies the name of the account to create. This is specified as part of the URL.

#### Sample Response

The example below shows output for the successful creation of `/ethereum/wallets/wallet1/accounts/account1/deposit-data`.

```
{
    "request_id": "b767dcca-5b10-4a52-1d9a-0a9b81b378ae",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "amount": 32000000000,
        "depositDataRoot": "6f3c48737eb58f47acf8268dad48d951577beb1d3c5dedbbe57f5ebf9e2c08a6",
        "publicKey": "g8yu7YB3jEy79oUubQf0ESg+WAfG9Vxh1zeUQr+3+wQhPHBL5Y3PRzma9CxQ5Zbc",
        "signature": "h2MWkBLgd4e6npSLdqWqJJI0opyPskph+IOIn6m4PS13rHOBaRHXaWarlm+0FymvF0OAWxPTMAa7jWiluOZJSryt3d+7Gqewdk15CGS2F2MCSAyUFCO4Kq6BR3kzvwyN",
        "withdrawalCredentials": "AK7JF5gDHMPr7gCzRT6bweOqw0ryIZ6VAasAlnAMisc="
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}
```