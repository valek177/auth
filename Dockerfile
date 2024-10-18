FROM golang:1.23.2-alpine3.20 AS builder

COPY . /github.com/valek177/auth/source/
WORKDIR /github.com/valek177/auth/source/

RUN go mod download
RUN go build -o ./bin/auth_server cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/valek177/auth/source/bin/auth_server .

CMD ["./auth_server"]
