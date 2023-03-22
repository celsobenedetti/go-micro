# base go image, used to build executable
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontApp ./cmd/web

RUN chmod +x /app/frontApp

FROM alpine:latest

COPY --from=builder /app/frontApp /app/frontApp

ENTRYPOINT [ "/app/frontApp" ]
