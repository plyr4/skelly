up: run-build-docker

run-build-docker: linux run-docker

clean:
	#################################
	######      Go clean       ######
	#################################

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "cleaning up"

build:
	#################################
	######    Build Binary     ######
	#################################
	@echo
	@echo "### Building static release/skelly binary"
	go build -o release/skelly github.com/davidvader/skelly/cmd/skelly

.PHONY: build-static-ci
build-static-ci:
	#################################
	######    Build CI Binary     ######
	#################################
	@echo
	@echo "### Building CI static release/skelly binary"
	@go build -a \
		-ldflags '-s -w -extldflags "-static" ${LD_FLAGS}' \
		-o release/skelly \
		github.com/davidvader/skelly/cmd/skelly

linux:
	#################################
	######  Build Linux Binary ######
	#################################
	@echo
	@echo "### Building static release/vela-server binary for linux"
	GOOS=linux CGO_ENABLED=0 \
		go build -o release/skelly github.com/davidvader/skelly/cmd/skelly

docker:
	#################################
	######    Build Docker     ######
	#################################
	@echo "### Building Docker image"
	docker build .

run:
	#################################
	######      Run Skelly     ######
	#################################
	@echo
	@echo "### Running skelly server"
	./release/skelly server

run-docker:
	#################################
	######    Restart Skelly   ######
	#################################
	@echo
	@echo "### Rebuilding and running skelly server"
	docker-compose up --build -d
