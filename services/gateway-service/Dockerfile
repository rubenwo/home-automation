FROM golang:1.17 AS builder

WORKDIR /go/gateway-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service ./cmd/api

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/gateway-service/gateway-service .

CMD [ "./gateway-service" ]