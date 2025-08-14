FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /build

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../../cmd ./cmd
COPY ../../internal ./internal
COPY ../../pkg ./pkg
COPY ../../service.dev.yaml ./

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -a -v -o api-gateway ./cmd/gateway/main.go

# ---- Final image ----
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /build/api-gateway ./
COPY --from=builder /build/service.dev.yaml ./

EXPOSE 8081

ENTRYPOINT ["./api-gateway"]
