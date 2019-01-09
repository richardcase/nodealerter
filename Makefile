
.PHONY: build
build: ## Build the controller
	@go build ./cmd/nodealerter

.PHONY: release
release: ## Build a release version of controller
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s" -o $(GOPATH)/bin/nodealerter ./cmd/nodealerter

.PHONY: docker-build
docker-build: ## Build a release version of controller via multi-stage docker file
	docker build .
