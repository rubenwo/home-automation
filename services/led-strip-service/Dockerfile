FROM golang:1.15 AS builder

WORKDIR /go/led-strip-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o led-strip-service ./cmd/api

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/led-strip-service/led-strip-service .

CMD [ "./led-strip-service" ]