#!/bin/bash

mkdir /data
vault server -config=/vault/config/vault-config.json -log-level=debug > /data/logs 2&1
sleep 5
if [ "$UNSEAL" = "true" ]; then
  /bin/sh /vault/config/vault-init.sh
  /bin/sh /vault/config/vault-unseal.sh
  /bin/sh /vault/config/vault-plugin.sh

fi

sleep 356000d
