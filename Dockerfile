FROM golang:1.22.1-alpine AS builder

ARG VER=dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --update make git gcc libc-dev curl file

RUN --mount=type=cache,target=/root/.cache/go-build \
    export CGO_ENABLED=1 &&\
    export LDFLAGS="-X main.version=$VER -linkmode external -extldflags \"-static\"" &&\
    go build -v -ldflags "${LDFLAGS}" -o shortener ./cmd/shortener/main.go &&\
    file shortener

FROM alpine:3.20.0

RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

RUN adduser -D shortener
USER shortener

WORKDIR /app

RUN mkdir /app/data

COPY --from=builder /app/shortener .

CMD ["./shortener"]
