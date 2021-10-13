FROM golang:1.17 AS builder

WORKDIR /go/tradfri-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tradfri-service ./cmd/api

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/tradfri-service/tradfri-service .
COPY --from=builder /go/tradfri-service/database ./database

CMD [ "./tradfri-service" ]