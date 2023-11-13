DB_URL=postgresql://gocommmicro:password@localhost:5432/gocommmicro?sslmode=disable
ADMIN_BACK_END_BINARY=adminBackEndApp
BROKER_BINARY=brokerApp
AUTHENTICATION_BINARY=authenticationApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the authentication binary as a linux executable
build_auth:
	@echo "Building authentication binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTHENTICATION_BINARY} ./cmd/api
	@echo "Done!"


## build_admin_back: builds the frone end binary
build_admin:
	@echo "Building admin backend binary..."
	cd ../admin-backend && env CGO_ENABLED=0 go build -o ${ADMIN_BACK_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the admin backend
start: build_admin
	@echo "Starting admin backend"
	cd ../admin-backend && ./${ADMIN_BACK_END_BINARY} &

## stop: stop the admin backend
stop:
	@echo "Stopping admin backend..."
	@-pkill -SIGTERM -f "./${ADMIN_BACK_END_BINARY}"
	@echo "Stopped admin backend!"

migrateup:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration/ -database "$(DB_URL)" -verbose down