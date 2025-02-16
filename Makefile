TAG = $(shell git describe --tags --always --dirty)
APP_NAME = $(shell basename $(CURDIR))
DOCKER_BASE = "ghcr.io/jamesmcdonald/$(APP_NAME)"
DOCKER_IMAGE = $(DOCKER_BASE):$(TAG)

$(APP_NAME): cmd/$(APP_NAME)/*.go cmd/$(APP_NAME)/version.go
	go build ./cmd/$(APP_NAME)

cmd/$(APP_NAME)/version.go: cmd/$(APP_NAME)/genversion.go
	go generate ./...

test: cmd/$(APP_NAME)/version.go
	go test ./...

docker:
	docker build -t $(DOCKER_BASE) .
	docker build -t $(DOCKER_IMAGE) .

push: docker
	docker push $(DOCKER_BASE)
	docker push $(DOCKER_IMAGE)

clean:
	docker rmi -f $(DOCKER_IMAGE)
	rm -f cmd/$(APP_NAME)/version.go $(APP_NAME) $(APP_NAME)-*.tgz

.PHONY: docker push clean generate build
