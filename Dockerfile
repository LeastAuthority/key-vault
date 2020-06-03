FROM vault

RUN apk add git
RUN mkdir /data/ && \
    cd /data/ && \
    git clone https://github.com/bloxapp/vault-plugin-secrets-eth2.0
