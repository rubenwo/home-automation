FROM golang:1.17 AS builder

WORKDIR /go/actions-service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o actions-service ./cmd/actions

FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/actions-service/actions-service .

CMD [ "./actions-service" ]