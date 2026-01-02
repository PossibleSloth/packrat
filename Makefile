# --- Variables ---
DOCKERHUB_USER = sl0th
IMAGE_NAME = packrat
GIT_HASH := $(shell git rev-parse --short HEAD)

# Full image name with tag
FULL_IMAGE_NAME = $(DOCKERHUB_USER)/$(IMAGE_NAME):$(GIT_HASH)

# --- Targets ---

# Default target runs the build command
all: build

# Target to build the Docker image
build:
	@echo "--- Building Docker image: $(FULL_IMAGE_NAME) ---"
	docker buildx build --platform linux/amd64 -t $(FULL_IMAGE_NAME) .
	@echo "--- Build complete ---"

# Target to push the Docker image to Docker Hub
push: build
	@echo "--- Pushing Docker image: $(FULL_IMAGE_NAME) to Docker Hub ---"
	docker push $(FULL_IMAGE_NAME)
	@echo "--- Push complete ---"

# Target to tag current version as "latest" tag
release:
	docker pull --platform linux/amd64 $(FULL_IMAGE_NAME)
	docker tag $(FULL_IMAGE_NAME) $(DOCKERHUB_USER)/$(IMAGE_NAME):latest
	docker push $(DOCKERHUB_USER)/$(IMAGE_NAME):latest

# Target to clean up local images (optional)
clean:
	@echo "--- Removing local Docker image: $(FULL_IMAGE_NAME) ---"
	docker rmi $(FULL_IMAGE_NAME)
	@echo "--- Removal complete ---"

# Target to run tests
test:
	go test ./...

# Run locally for testing
run:
	LISTEN_HOST=127.0.0.1 LISTEN_PORT=7777 STATIC_DIR=/tmp SERVER_HOST=127.0.0.1:7777 go run main.go

.PHONY: all build push test release run clean
