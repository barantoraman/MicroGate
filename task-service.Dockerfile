FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

RUN apk add git && \
    git clone https://github.com/grpc-ecosystem/grpc-health-probe.git

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./service.dev.yaml  ./

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -a -v -o task ./cmd/task/main.go

WORKDIR /build/grpc-health-probe
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o grpc-health-probe .
WORKDIR /build

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /build/task ./
COPY --from=builder /build/service.dev.yaml ./
COPY --from=builder /build/grpc-health-probe/grpc-health-probe ./

EXPOSE 8082
ENTRYPOINT ./task
