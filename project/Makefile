DB_BINARY=DBApp
BACK_END_BINARY=backendApp
## up: starts all containers in the background without forcing build
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build : stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_db 
	@echo "Stopping docker images(if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## up_down: stop docker compose
up_down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_db: builds the db binary as a linux executable
build_db:
	@echo "Building DB binary..."
	cd ../db && protoc -Iproto --go_opt=module=github.com/GirishBhutiya/gocomm-micro/db --go_out=. --go-grpc_opt=module=github.com/GirishBhutiya/gocomm-micro/db --go-grpc_out=. proto/*.proto
	cd ../db && env GOOS=linux CGO_ENABLED=0 go build -o ${DB_BINARY} ./cmd/api
	@echo "Done!"

## build_back: builds the back end binary
build_back:
	@echo Building back end binary...
	chdir ../backend && env CGO_ENABLED=0 && env GOOS=linux && go build -o ${BACK_END_BINARY} ./cmd/web
	@echo Done!

## start: starts the back end
start: build_back
	@echo Starting back end
	chdir ../backend && ./${BACK_END_BINARY} &

## stop: stop the back end
stop:
	@echo Stopping back end...
	@-pkill -SIGTERM -f ${BACK_END_BINARY}
	@echo "Stopped back end!"