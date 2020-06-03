FROM vault


RUN mkdir /data/ && \
    cd /data/src && \
    git clone https://github.com/bloxapp/vault-plugin-secrets-eth2.0
