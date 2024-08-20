FROM golang:1.23-alpine AS builder

ARG CURRENT_TIME
ARG GIT_DESCRIPTION

ENV buildTime=${CURRENT_TIME}
ENV version=${GIT_DESCRIPTION}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags "-s -X main.buildTime=$buildTime -X main.version=$version" -o ./bin/api ./cmd/api/

FROM alpine:3.18

ENV MIGRATE_VERSION=v4.17.1
ENV OS=linux
ENV ARCH=amd64

RUN apk --no-cache add curl

# Download and install migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.${OS}-${ARCH}.tar.gz \
    | tar xz -C /usr/local/bin

COPY --from=builder /app/bin/api /usr/local/bin/api

EXPOSE 4000

ARG DB_DSN
ENV dbDsn=${DB_DSN}

COPY ./scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /migrations-dir
COPY ./migrations .

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]