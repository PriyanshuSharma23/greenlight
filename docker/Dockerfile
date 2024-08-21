FROM golang:1.23-alpine AS builder

ARG CURRENT_TIME
ARG GIT_DESCRIPTION

ENV buildTime=${CURRENT_TIME}
ENV version=${GIT_DESCRIPTION}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -ldflags "-s -X main.buildTime=$buildTime -X main.version=$version" -o ./bin/api ./cmd/api/

FROM alpine:3.18

COPY --from=builder /app/bin/api /usr/local/bin/api

ARG DB_DSN
ENV dbDsn ${DB_DSN}

EXPOSE 4000
ENTRYPOINT /usr/local/bin/api -db-dsn=$dbDsn