FROM golang:1.15 AS builder

WORKDIR /go/hue-service
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o hue-service ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/hue-service/hue-service .

CMD [ "./hue-service" ]