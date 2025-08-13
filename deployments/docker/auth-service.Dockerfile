FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /build

COPY ../../go.mod ../../go.sum ./
RUN go mod download

RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

RUN apk add --no-cache git && \
    git clone https://github.com/grpc-ecosystem/grpc-health-probe.git

COPY ../../cmd ./cmd
COPY ../../internal ./internal
COPY ../../pkg ./pkg
COPY ../../service.dev.yaml ./

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -a -v -o auth ./cmd/auth/main.go

RUN cd ./grpc-health-probe && \
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o grpc-health-probe .

# ---- Final image ----
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /build/auth ./
COPY --from=builder /build/internal/auth/db/migrations ./auth-migrations
COPY --from=builder /build/service.dev.yaml ./
COPY --from=builder /build/migrate ./
COPY --from=builder /build/grpc-health-probe/grpc-health-probe ./

EXPOSE 8083

CMD ./migrate -path=./auth-migrations -database="postgres://microgate:example@postgres:5432/auth?sslmode=disable" up && \
    ./auth
