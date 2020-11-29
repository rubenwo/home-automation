FROM golang:1.15 AS builder

WORKDIR /go/gateway-service
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway-service ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/gateway-service/gateway-service .
#COPY --from=builder /go/gateway-service/ingress.yaml .
#COPY --from=builder /go/gateway-service/database.json .
#
#COPY --from=builder /go/gateway-service/assets/crypto/ ./assets/crypto/

CMD [ "./gateway-service" ]