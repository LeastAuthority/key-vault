# vault-plugin-secrets-eth2.0
[![blox.io](https://s3.us-east-2.amazonaws.com/app-files.blox.io/static/media/powered_by.png)](https://blox.io)

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
