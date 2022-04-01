FROM golang:1.17 AS builder

WORKDIR /go/qr-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o qr-service ./cmd/app

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/qr-service/qr-service .

CMD [ "./qr-service" ]