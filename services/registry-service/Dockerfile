FROM golang:1.16 AS builder

WORKDIR /go/registry-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o registry-service ./cmd/api

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/registry-service/registry-service .
COPY --from=builder /go/registry-service/config.json .

CMD [ "./registry-service" ]