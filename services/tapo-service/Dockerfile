FROM golang:1.20 AS builder

WORKDIR /go/tapo-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tapo-service ./cmd/api

FROM alpine:3.12 AS alpine
RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /
COPY --from=builder /go/tapo-service/tapo-service .
COPY --from=builder /go/tapo-service/config.json .
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "./tapo-service" ]