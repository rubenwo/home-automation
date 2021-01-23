FROM golang:1.15 AS builder

WORKDIR /go/hue-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o hue-service ./cmd/api

FROM alpine:3.7
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/hue-service/hue-service .

CMD [ "./hue-service" ]