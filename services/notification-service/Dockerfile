FROM golang:1.17 AS builder

WORKDIR /go/notification-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/app

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/notification-service/notification-service .
COPY --from=builder /go/notification-service/config.json .

CMD [ "./notification-service" ]