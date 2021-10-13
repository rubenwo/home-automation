FROM golang:1.17 AS builder

WORKDIR /go/video-streaming-hub-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o video-streaming-hub-service ./cmd/app

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/video-streaming-hub-service/video-streaming-hub-service .

CMD [ "./video-streaming-hub-service" ]