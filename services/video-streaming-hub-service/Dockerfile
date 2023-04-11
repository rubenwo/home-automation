FROM golang:1.20 AS builder

WORKDIR /go/video-streaming-hub-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o video-streaming-hub-service ./cmd/video-streaming-hub

FROM alpine:3.12 AS alpine
RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /
COPY --from=builder /go/video-streaming-hub-service/video-streaming-hub-service .
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD [ "./video-streaming-hub-service" ]