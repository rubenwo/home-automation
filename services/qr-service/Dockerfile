FROM golang:1.20 AS builder

WORKDIR /go/qr-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o qr-service ./cmd/app

FROM alpine:3.12 AS alpine
RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /
COPY --from=builder /go/qr-service/qr-service .
COPY --from=builder /go/qr-service/config.json .
COPY --from=builder /go/qr-service/logo.png .
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "./qr-service" ]