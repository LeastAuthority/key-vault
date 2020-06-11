# vault-plugin-secrets-eth2.0

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
        "wallets": {
            "{\"crypto\":{\"original\":\"R/DQm1NoDtbvZQm3OK+cb7A4qipMEqbdwgF5kwy/4lE=\"},\"name\":\"wallet1\",\"nextaccount\":3,\"type\":\"hierarchical deterministic\",\"uuid\":\"80eb953c-4478-42bb-add3-7b4e143d4051\",\"version\":1}": true
        }
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
        "id": "f461caab-0bd7-4f81-9816-f055dc027a81",
        "name": "wallet2"
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
        "accounts": {
            "{\"crypto\":{\"original\":\"EGtt6Ujj+kYGfwOkUCPqcfKFhxVMskpbOBtCi5XtW1A=\"},\"name\":\"account1\",\"path\":\"m/12381/3600/0/0\",\"pubkey\":\"b9c4c809fb2536b76107bea7fde58bd686e5b593358e163868bdf4a006d6128a459f7fa62064a6ef3d94c804ae6efa6a\",\"uuid\":\"f164c78b-1ba7-4f59-9d97-2e129b12a992\",\"version\":1}": true,
            "{\"crypto\":{\"original\":\"JcjxrkBnvinsKqW6ycZF6J3ToiopgWpYPD51dXYVl6c=\"},\"name\":\"account1\",\"path\":\"m/12381/3600/1/0\",\"pubkey\":\"a8fb2d7fa997b3db9a004866ba8b52ec4534213c969d3813197bc2e590b54146f6033eea3d46cf97bd291bb479535320\",\"uuid\":\"8e92540a-8403-4329-911e-af5f8afea7c9\",\"version\":1}": true
        }
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

The example below shows output for the successful creation of `/ethereum/wallets/wallet1/account1`.

```
{
    "request_id": "b767dcca-5b10-4a52-1d9a-0a9b81b378ae",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "accountName": "account1",
        "path": "m/12381/3600/2/0",
        "walletName": "wallet1"
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
| `POST`  | `:mount-path/wallets/:wallet_name/accounts/:account_name/sign`  | `200 application/json` |

#### Parameters

* `walet_name` (`string: <required>`) - Specifies the name of the wallet of the account to sign. This is specified as part of the URL.
* `account_name` (`string: <required>`) - Specifies the name of the account to sign. This is specified as part of the URL.

#### Sample Response

The example below shows output for the successful creation of `/ethereum/wallets/wallet1/account1/sign`.

```
{
    "request_id": "b767dcca-5b10-4a52-1d9a-0a9b81b378ae",
    "lease_id": "",
    "renewable": false,
    "lease_duration": 0,
    "data": {
        "accountName": "account1",
        "path": "m/12381/3600/2/0",
        "walletName": "wallet1"
    },
    "wrap_info": null,
    "warnings": null,
    "auth": null
}
```