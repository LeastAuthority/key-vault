FROM vault

RUN apk add git
RUN apk add go
RUN apk add --update gcc g++

COPY ./config/vault-config.json /vault/config/vault-config.json

RUN pwd

RUN ls -lah

COPY ./backend /backend
COPY go.mod /
COPY go.sum /
COPY main.go /

RUN go build -o ethsign main.go && \
    mkdir /vault/plugins/ && \
    cp ethsign /vault/plugins/
