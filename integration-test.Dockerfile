FROM golang:1.24.5-alpine3.22

WORKDIR /tests

COPY go.mod go.sum ./
RUN go mod download

COPY ./integration_tests ./

CMD ["go", "test", "-v", "-count=1"]
