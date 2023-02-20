# minimal Docker image simply to run binary created with Makefile
FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp" ]
