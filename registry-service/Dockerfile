FROM golang:1.15 AS builder

WORKDIR /go/registry-service
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o registry-service ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/registry-service/registry-service .

CMD [ "./registry-service" ]