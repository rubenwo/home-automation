FROM golang:1.20 AS builder

WORKDIR /go/led-strip-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o led-strip-service ./cmd/api

FROM alpine:3.12 AS alpine
RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /
COPY --from=builder /go/led-strip-service/led-strip-service .
COPY --from=builder /go/led-strip-service/config.json .
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "./led-strip-service" ]