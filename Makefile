# Variables
APP_NAME = proxy-app
GO_CMD = go
DOCKER_CMD = docker
DOCKER_IMAGE = ashok-proxy-image

run:
	@$(GO_CMD) run cmd/proxy/main.go

build:
	@$(GO_CMD) build -o $(APP_NAME) cmd/proxy/main.go

test:
	@$(GO_CMD) test ./...

docker-build:
	@$(DOCKER_CMD) build -t $(DOCKER_IMAGE) .

docker-run:
	@$(DOCKER_CMD) run -p 8080:8080 $(DOCKER_IMAGE)
setc:
	@python scripts/cookie-getter.py
air:
	@air