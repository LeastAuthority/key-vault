#!/bin/sh

vault server -config=/vault/config/vault-config.json -log-level=debug > /data/logs 2&1
sleep 5
if [ "$UNSEAL" = "true" ]; then

  /bin/sh /vault/config/vault-init.sh
  sleep 5
  /bin/sh /vault/config/vault-unseal.sh
  sleep 5
  /bin/sh /vault/config/vault-plugin.sh
  sleep 5

fi

sleep 356000d
