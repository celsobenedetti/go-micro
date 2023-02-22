.DEFAULT_GOAL := help

FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp

##
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

UP_DESC="up: starts all containers in the background without forcing build"

##
up_build: build_broker build_auth build_logger
	@echo "Stopping any running Docker containers"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

UPBUILD_DESC="up_build: stops docker-compose \(if running\), builds all projects and starts docker compose"

##
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

DOWN_DESC="down: stop docker compose"

##
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

BUILD_BROKER_DESC="build_broker: builds the broker binary as a linux executable"

##
build_auth:
	@echo "Building auth service binary..."
	cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

BUILD_AUTH_DESC="build_auth: builds the auth service binary as a linux executable"

##
build_logger:
	@echo "Building logger service binary..."
	cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

BUILD_logger_DESC="build_logger: builds the logger service binary as a linux executable"

##
build_front:
	@echo "Building front end binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

BUILD_FRONT_DESC="build_front: builds the frone end binary"

##
start: build_front
	@echo "Starting front end"
	cd ./front-end && ./${FRONT_END_BINARY} &

START_DESC="start: starts the front end"

##
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

STOP_DESC="stop: stop the front end"

##
help:
	@echo "Usage:\n\tmake <target> \n\nTargets:"
	@echo "\t${UP_DESC}"
	@echo "\t${UPBUILD_DESC}"
	@echo "\t${DOWN_DESC}"
	@echo "\t${BUILD_BROKER_DESC}"
	@echo "\t${BUILD_AUTH_DESC}"
	@echo "\t${BUILD_FRONT_DESC}"
	@echo "\t${START_DESC}"
	@echo "\t${STOP_DESC}"
