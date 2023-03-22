FROM alpine:latest

COPY ./frontApp /app/frontApp

RUN chmod +x /app/frontApp

ENTRYPOINT [ "/app/frontApp" ]
