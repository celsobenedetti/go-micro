# base go image, used to build executable
# not needed since makefile can create binary if system has go installed

# FROM golang:1.18-alpine as builder
#
# RUN mkdir /app
#
# COPY . /app
#
# WORKDIR /app
#
# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api
#
# RUN chmod +x /app/brokerApp


# actual app running image
# build a tiny docker image, copy binary and run
FROM alpine:latest

RUN mkdir /app

# example: copy binary build by 'builder' image
# COPY --from=builder /app/brokerApp /app 

# copy binary build by Makefile (much faster)
COPY brokerApp /app 

CMD [ "/app/brokerApp" ]
