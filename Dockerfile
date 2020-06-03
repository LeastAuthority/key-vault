FROM vault

RUN apk add go

COPY ./config/vault-config.json /vault/config/vault-config.json

RUN go build -o ethsign main.go && \
    mkdir /vault/plugins/ && \
    cp ethsign /vault/plugins/
