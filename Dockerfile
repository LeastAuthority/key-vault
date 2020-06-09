FROM vault

RUN apk add git
RUN apk add go
RUN apk add --update gcc g++

# enable Go modules support
ENV GO111MODULE=on

COPY go.mod /
COPY go.sum /
RUN go mod download

COPY ./config/vault-config.json /vault/config/vault-config.json
COPY ./config/vault-init-unseal.sh /vault/config/vault-init-unseal.sh
COPY ./backend /backend
COPY main.go /

RUN CGO_ENABLED=1 GOOS=linux go build -a -o ethsign main.go
RUN mkdir /vault/plugins/ && \
    cp ethsign /vault/plugins/
